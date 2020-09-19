package ffmpegLib

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type IgnoreErrorHandler struct {
	Handler Handler
}

func (h *IgnoreErrorHandler) Progress() float32 {
	return h.Handler.Progress()
}

func (h *IgnoreErrorHandler) Run(progressHandler ProgressListener) error {
	logrus.Infoln(fmt.Sprintf("Starting %s", reflect.TypeOf(h).String()))
	_ = h.Handler.Run(progressHandler)
	return nil
}

func NewIgnoreErrorHandler(handler Handler) *IgnoreErrorHandler {
	return &IgnoreErrorHandler{
		Handler: handler,
	}
}
