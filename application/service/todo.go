package service

import (
	"github.com/schidstorm/ffmpeg-go-server/application/model"
	"gorm.io/gorm"
)

type Todo struct {
	model.Todo
}

func (u *Todo) AfterCreate(tx *gorm.DB) (err error) {

	return nil
}
