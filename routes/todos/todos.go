package todos

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gdb "github.com/nickbullock/pomodoro-api/db"
	todo "github.com/nickbullock/pomodoro-api/models/todo"
)

var db *gorm.DB
var v1 *gin.RouterGroup

type TodoModel = todo.Model

func Register(r *gin.Engine) {
	db = gdb.GetDB()
	v1 = r.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
}

// createTodo add a new todo
func createTodo(c *gin.Context) {
	var todo = TodoModel{}

	c.BindJSON(&todo)

	db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo item created successfully", "resourceId": todo.ID})
}

// fetchAllTodo fetch all todos
func fetchAllTodo(c *gin.Context) {
	var todos []TodoModel

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todos})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todos})
}

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo TodoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todo})
}

// updateTodo update a todo
func updateTodo(c *gin.Context) {
	var todo = TodoModel{}
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found"})
		return
	}

	c.BindJSON(&todo)

	db.Model(&todo).Update("title", todo.Title)
	db.Model(&todo).Update("status", todo.Status)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully"})
}

// deleteTodo remove a todo
func deleteTodo(c *gin.Context) {
	var todo TodoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully"})
}
