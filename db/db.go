package db

import (
	"flag"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	group "github.com/nickbullock/pomodoro-api/models/group"
	todo "github.com/nickbullock/pomodoro-api/models/todo"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Init() {
	//open a db connection
	var err error

	var envMode = os.Getenv("MODE")
	var envFile string = ".env"

	flag.Parse()

	if envMode == "PROD" {
		envFile = ".env-prod"
	}

	fmt.Println("mode:", envMode, envFile)

	err = godotenv.Load(envFile)
	if err != nil {
		panic("failed to load env")
	}

	pass := os.Getenv("DB_PASS")

	fmt.Println("db:", os.Getenv("DB_NAME"), pass)

	dbURLemplate := "host=%s port=%s user=%s dbname=%s sslmode=disable"
	dbOptions := []interface{}{os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME")}
	if pass != "" {
		dbURLemplate += " password=%s"
		dbOptions = append(dbOptions, pass)
	}
	dbURL := fmt.Sprintf(dbURLemplate, dbOptions...)

	fmt.Println("db url:", dbURL)

	db, err = gorm.Open("postgres", dbURL)
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	db.AutoMigrate(&group.Model{})
	db.AutoMigrate(&todo.Model{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
}
