package ffmpegLib

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
)

type RemoveHandler struct {
	FilePath string
}

func (h *RemoveHandler) Progress() float32 {
	return 0
}

func (h *RemoveHandler) Run(progressHandler ProgressListener) error {
	logrus.Infoln(fmt.Sprintf("Starting %s", reflect.TypeOf(h).Name()))
	return os.Remove(h.FilePath)
}

func NewRemoveHandler(filePath string) *RemoveHandler {
	return &RemoveHandler{
		FilePath: filePath,
	}
}
