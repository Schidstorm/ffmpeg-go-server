package ffmpegLib

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type MP4TranscoderHandler struct {
	SourceFilePath      string
	DestinationFilePath string
	FfmpegMultiHandler  Handler
}

func (h *MP4TranscoderHandler) Progress() float32 {
	return h.FfmpegMultiHandler.Progress()
}

func (h *MP4TranscoderHandler) Run(progressHandler ProgressListener) error {
	logrus.Infoln(fmt.Sprintf("Starting %s", reflect.TypeOf(h).String()))
	return h.FfmpegMultiHandler.Run(progressHandler)
}

func NewMP4TranscoderHandler(sourceDirectory, destinationDirectory string) *MP4TranscoderHandler {
	return &MP4TranscoderHandler{
		SourceFilePath:      sourceDirectory,
		DestinationFilePath: destinationDirectory,
		FfmpegMultiHandler: NewMultiHandler(
			NewFfmpegHandler("-y", "-i", sourceDirectory, "-acodec", "aac", "-c:v", "libx264", "-crf", "18", "-preset", "veryslow", destinationDirectory),
		),
	}
}
