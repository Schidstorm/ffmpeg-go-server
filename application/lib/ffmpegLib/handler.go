package ffmpegLib

type ProgressCallback = func(progress float32)
type FinishCallback func(err error)

type Handler interface {
	Progress() float32
	Run(listener ProgressListener) error
}
