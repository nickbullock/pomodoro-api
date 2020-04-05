package main

import (
	"os"

	db "github.com/nickbullock/pomodoro-api/db"
	groups "github.com/nickbullock/pomodoro-api/routes/groups"
	todos "github.com/nickbullock/pomodoro-api/routes/todos"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	db.Init()
}

func main() {
	var port = os.Getenv("PORT")
	router := gin.Default()
	router.Use(cors.Default())

	todos.Register(router)
	groups.Register(router)

	router.Run(":" + port)
}
