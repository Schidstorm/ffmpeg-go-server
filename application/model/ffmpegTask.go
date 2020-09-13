package model

import (
	"gorm.io/gorm"
	"time"
)

type FfmpegTask struct {
	gorm.Model

	SourceFileName      string
	DestinationFileName string
	Started             bool
	StartedAt           time.Time
	Finished            bool
	FinishedAt          time.Time
	Progress            float32
}
