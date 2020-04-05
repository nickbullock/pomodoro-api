package group

import (
	pq "github.com/lib/pq"
	gorm "github.com/nickbullock/pomodoro-api/models/gorm"
	todo "github.com/nickbullock/pomodoro-api/models/todo"
)

type (
	// Model describes a Model type
	Model struct {
		gorm.Model
		Title   string         `json:"title"`
		Creator string         `json:"creator"`
		Members pq.StringArray `json:"members" gorm:"type:varchar(100)[]"`
		// association?
		Todos []todo.Model `json:"todos" gorm:"foreignkey:GroupID"`
	}
)

func (group *Model) TableName() string {
	return "groups"
}
