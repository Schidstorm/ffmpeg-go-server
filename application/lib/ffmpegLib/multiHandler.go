package ffmpegLib

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type MultiHandler struct {
	handlers            []Handler
	currentHandlerIndex int
}

func (h *MultiHandler) Progress() float32 {
	divider := float32(1) / float32(len(h.handlers))
	return float32(h.currentHandlerIndex)*divider + h.handlers[h.currentHandlerIndex].Progress()/divider
}

func (h *MultiHandler) Run(progressHandler ProgressListener) error {
	logrus.Infoln(fmt.Sprintf("Starting %s", reflect.TypeOf(h).Name()))
	for index, handler := range h.handlers {
		h.currentHandlerIndex = index
		err := handler.Run(progressHandler)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewMultiHandler(handlers ...Handler) *MultiHandler {
	return &MultiHandler{
		handlers:            handlers,
		currentHandlerIndex: 0,
	}
}
