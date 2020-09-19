package ffmpegLib

import (
	"fmt"
	"github.com/schidstorm/ffmpeg-go-server/application/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

	sourceFilePath := path.Join(GetSettings().UploadDestinationDirectory, model.SourceFileName)
	tempFilePath := path.Join(GetSettings().TempDirectory, model.DestinationFileName)
	destinationFilePath := path.Join(GetSettings().ConversionDestinationDirectory, model.DestinationFileName)

	handler := MultiHandler{
		handlers: []Handler{
			NewMP4TranscoderHandler(sourceFilePath, tempFilePath),
			NewIgnoreErrorHandler(NewRemoveHandler(destinationFilePath)),
			NewCopyHandler(tempFilePath, destinationFilePath),
			NewRemoveHandler(tempFilePath),
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
