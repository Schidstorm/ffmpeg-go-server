package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model

	Name    string `json:"name"`
	Message string `json:"message"`
}
