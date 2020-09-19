package ffmpegLib

type MP4TranscoderHandler struct {
	SourceFilePath      string
	DestinationFilePath string
	FfmpegMultiHandler  Handler
}

func (h *MP4TranscoderHandler) Progress() float32 {
	return h.FfmpegMultiHandler.Progress()
}

func (h *MP4TranscoderHandler) Run(progressHandler ProgressListener) error {
	return h.FfmpegMultiHandler.Run(progressHandler)
}

func NewMP4TranscoderHandler(sourceDirectory, destinationDirectory string) *MP4TranscoderHandler {
	return &MP4TranscoderHandler{
		SourceFilePath:      sourceDirectory,
		DestinationFilePath: destinationDirectory,
		FfmpegMultiHandler: NewMultiHandler(
			NewFfmpegHandler("-y", "-i", sourceDirectory, "-vcodec", "h264", "-acodec", "aac", "-strict", "-2", destinationDirectory),
		),
	}
}
