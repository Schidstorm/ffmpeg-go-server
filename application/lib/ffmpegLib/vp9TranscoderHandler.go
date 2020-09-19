package ffmpegLib

import (
	"os"
)

type VP9TranscoderHandler struct {
	SourceFilePath      string
	DestinationFilePath string
	FfmpegMultiHandler  Handler
}

func (h *VP9TranscoderHandler) Progress() float32 {
	return h.FfmpegMultiHandler.Progress()
}

func (h *VP9TranscoderHandler) Run(progressHandler ProgressListener) error {
	return h.FfmpegMultiHandler.Run(progressHandler)
}

func NewVP9TranscoderHandler(sourceDirectory, destinationDirectory string) *VP9TranscoderHandler {
	return &VP9TranscoderHandler{
		SourceFilePath:      sourceDirectory,
		DestinationFilePath: destinationDirectory,
		FfmpegMultiHandler: NewMultiHandler(
			NewFfmpegHandler(
				"-y", "-re", "-hide_banner", "-progress", "pipe:2", "-i", sourceDirectory,
				"-c:v", "libvpx-vp9", "-pass", "1", "-b:v", "1000K", "-threads", "6", "-speed", "4", "-tile-columns", "0", "-frame-parallel", "0", "-auto-alt-ref", "1", "-lag-in-frames", "25", "-g", "9999", "-aq-mode", "0", "-an", "-f", "webm", os.DevNull),
			NewFfmpegHandler(
				"-y", "-re", "-hide_banner", "-progress", "pipe:2", "-i", sourceDirectory,
				"-c:v", "libvpx-vp9", "-pass", "2", "-b:v", "1000K", "-threads", "6", "-speed", "0", "-tile-columns", "0", "-frame-parallel", "0", "-auto-alt-ref", "1", "-lag-in-frames", "25", "-g", "9999", "-aq-mode", "0", "-c:a", "libopus", "-b:a", "64k", "-f", "webm", destinationDirectory),
		),
	}
}
