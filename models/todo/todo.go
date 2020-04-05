package todo

import (
	gorm "github.com/nickbullock/pomodoro-api/models/gorm"
)

type Status int

const (
	New Status = iota
	InProgress
	Done
)

func (d Status) String() string {
	return [...]string{"New", "InProgress", "Done"}[d]
}

type (
	// Model describes a Model type
	Model struct {
		gorm.Model
		Title      string `json:"title"`
		Status     string `json:"status"`
		GroupID    uint   `json:"groupId"`
		InProgress uint   `json:"inProgress"`
	}
)

func (todo *Model) TableName() string {
	return "todos"
}
