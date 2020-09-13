package service

import (
	"github.com/schidstorm/ffmpeg-go-server/application/model"
	"gorm.io/gorm"
)

type FfmpegTask struct {
	model.FfmpegTask
}

func (u *FfmpegTask) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", "admin")
	}

	return
}
