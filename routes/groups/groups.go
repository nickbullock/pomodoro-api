package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gdb "github.com/nickbullock/pomodoro-api/db"
	"github.com/nickbullock/pomodoro-api/models/group"
)

var db *gorm.DB
var v1 *gin.RouterGroup

func Register(r *gin.Engine) {
	db = gdb.GetDB()
	v1 = r.Group("/api/v1/groups")
	{
		v1.POST("/", createGroup)
		v1.GET("/", getGroups)
		v1.GET("/:id", getGroup)
		v1.PUT("/:id", updateGroup)
		v1.DELETE("/:id", deleteGroup)
	}
}

// createTodo add a new todo
func createGroup(c *gin.Context) {
	var group = group.Model{}
	c.BindJSON(&group)
	q := c.Request.URL.Query()

	if q["userId"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "No user id"})
		return
	}

	userID := q["userId"][0]

	group.Creator = userID
	group.Members = append(group.Members, userID)

	db.Save(&group)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Group created successfully!", "resourceId": group.ID})
}

func getGroups(c *gin.Context) {
	var groups []group.Model

	q := c.Request.URL.Query()

	if q["userId"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "No user id"})
		return
	}

	userID := q["userId"][0]

	db.Preload("Todos", func(db *gorm.DB) *gorm.DB {
		return db.Order("todos.created_at")
	}).Where("?=ANY(members) AND deleted_at IS NULL", userID).Order("created_at").Find(&groups)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": groups})
}

func deleteGroup(c *gin.Context) {
	var group group.Model
	groupID := c.Param("id")

	db.First(&group, groupID)

	if group.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No group found"})
		return
	}

	db.Delete(&group)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Group deleted successfully"})
}

func getGroup(c *gin.Context) {
	var group group.Model
	groupID := c.Param("id")

	q := c.Request.URL.Query()

	if q["userId"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "No user id"})
		return
	}

	userID := q["userId"][0]

	db.Preload("Todos").First(&group, groupID).Where("?=ANY(members) AND deleted_at IS NULL", userID)

	if group.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No group found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": group})
}

func updateGroup(c *gin.Context) {
	var group = group.Model{}
	groupID := c.Param("id")

	db.First(&group, groupID)

	if group.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No group found"})
		return
	}

	c.BindJSON(&group)

	if group.Members != nil {
		db.Model(&group).Update("members", group.Members)
	}
	if group.Title != "" {
		db.Model(&group).Update("title", group.Title)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Group updated successfully"})
}
