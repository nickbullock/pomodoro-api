package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	//open a db connection
	var err error

	var envMode = os.Getenv("MODE")
	var envFile string = ".env"

	flag.Parse()

	if envMode == "prod" {
		envFile = ".env-prod"
	}

	fmt.Println("mode:", envMode, envFile)

	err = godotenv.Load(envFile)
	if err != nil {
		panic("failed to load env")
	}

	fmt.Println("db:", os.Getenv("DB_NAME"))

	dbOptions := []interface{}{os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), "sslmode=disable"}
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s %s", dbOptions...)

	fmt.Println("db url:", dbURL, dbOptions)

	db, err = gorm.Open("postgres", dbURL)
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	db.AutoMigrate(&todoModel{})
}

func main() {
	var port = os.Getenv("PORT")
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}

	router.Run(":" + port)
	// router.RunTLS(":8080", "server.crt", "server.key")

}

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      gorm.Model
//    }
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

type (
	// todoModel describes a todoModel type
	todoModel struct {
		Model
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	// transformedTodo represents a formatted todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

// createTodo add a new todo
func createTodo(c *gin.Context) {
	var todo = todoModel{}
	c.BindJSON(&todo)

	db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

// fetchAllTodo fetch all todos
func fetchAllTodo(c *gin.Context) {
	var todos []todoModel

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todos})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todos})
}

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// updateTodo update a todo
func updateTodo(c *gin.Context) {
	var todo = todoModel{}
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	c.BindJSON(&todo)

	db.Model(&todo).Update("title", todo.Title)
	db.Model(&todo).Update("completed", todo.Completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// deleteTodo remove a todo
func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
