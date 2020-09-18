package ffmpegLib

import (
	"fmt"
	"github.com/schidstorm/ffmpeg-go-server/application/lib"
	"github.com/schidstorm/ffmpeg-go-server/application/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"path"
	"time"
)

type Consumer struct {
	db *gorm.DB
}

func NewConsumer(db *gorm.DB) *Consumer {
	return &Consumer{db: db}
}

func (consumer Consumer) Run() {
	go consumer.runAsync()
}

func (consumer Consumer) runAsync() {
	for {
		logrus.Infoln("Searching for ffmpeg tasks")
		model := &service.FfmpegTask{}
		model.Started = false

		if consumer.db.Where(map[string]interface{}{
			"finished": false,
		}).Take(model).Error == nil {
			consumer.runFFmpeg(model)
		}

		time.Sleep(10 * time.Second)
	}
}

func (consumer Consumer) runFFmpeg(model *service.FfmpegTask) {
	model.Started = true
	model.StartedAt = time.Now()
	consumer.db.Updates(model)

	sourceFilePath := path.Join(lib.GetSettings().UploadDestinationDirectory, model.SourceFileName)
	destinationFilePath := path.Join(lib.GetSettings().ConversionDestinationDirectory, model.DestinationFileName)

	handler := MultiHandler{
		handlers: []Handler{
			NewFfmpegHandler(
				"-y", "-re", "-hide_banner", "-progress", "pipe:2", "-i", sourceFilePath,
				"-c:v", "libvpx-vp9", "-pass", "1", "-b:v", "1000K", "-threads", "1", "-speed", "4", "-tile-columns", "0", "-frame-parallel", "0", "-auto-alt-ref", "1", "-lag-in-frames", "25", "-g", "9999", "-aq-mode", "0", "-an", "-f", "webm", os.DevNull),
			NewFfmpegHandler(
				"-y", "-re", "-hide_banner", "-progress", "pipe:2", "-i", sourceFilePath,
				"-c:v", "libvpx-vp9", "-pass", "2", "-b:v", "1000K", "-threads", "1", "-speed", "0", "-tile-columns", "0", "-frame-parallel", "0", "-auto-alt-ref", "1", "-lag-in-frames", "25", "-g", "9999", "-aq-mode", "0", "-c:a", "libopus", "-b:a", "64k", "-f", "webm", destinationFilePath),
		},
	}
	err := handler.Run(func(progress float32) {
		model.Progress = progress
		consumer.db.Updates(model)
		logrus.Infoln(progress)
	})
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Error converting %s", model.SourceFileName), err)
	} else {
		model.Finished = true
		model.FinishedAt = time.Now()
		consumer.db.Updates(model)
	}
}
