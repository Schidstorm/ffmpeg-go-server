package ffmpegLib

import (
	"os"
)

type RemoveHandler struct {
	FilePath string
}

func (h *RemoveHandler) Progress() float32 {
	return 0
}

func (h *RemoveHandler) Run(progressHandler ProgressListener) error {
	return os.Remove(h.FilePath)
}

func NewRemoveHandler(filePath string) *RemoveHandler {
	return &RemoveHandler{
		FilePath: filePath,
	}
}
