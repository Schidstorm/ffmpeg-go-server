package ffmpegLib

import (
	"io"
	"os"
)

type CopyHandler struct {
	SourceFilePath      string
	DestinationFilePath string
}

func (h *CopyHandler) Progress() float32 {
	return 0
}

func (h *CopyHandler) Run(progressHandler ProgressListener) error {
	source, err := os.OpenFile(h.SourceFilePath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	destination, err := os.OpenFile(h.DestinationFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024*1024*100)
	_, err = io.CopyBuffer(destination, source, buffer)
	if err != nil {
		return err
	}

	return nil
}

func NewCopyHandler(sourceDirectory, destinationDirectory string) *CopyHandler {
	return &CopyHandler{
		SourceFilePath:      sourceDirectory,
		DestinationFilePath: destinationDirectory,
	}
}
