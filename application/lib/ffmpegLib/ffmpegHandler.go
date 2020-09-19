package ffmpegLib

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type FfmpegHandler struct {
	arguments       []string
	totalDuration   time.Duration
	progressionTime time.Duration
}

func NewFfmpegHandler(arguments ...string) *FfmpegHandler {
	return &FfmpegHandler{
		arguments: arguments,
	}
}

func (h *FfmpegHandler) Progress() float32 {
	fDuration := float32(h.totalDuration)
	if fDuration == 0 {
		return 0
	}
	return float32(h.progressionTime) / float32(h.totalDuration)
}

func (h *FfmpegHandler) Run(handler ProgressListener) error {
	logrus.Infoln(fmt.Sprintf("Starting %s", reflect.TypeOf(h).Name()))
	cmd := exec.Command("ffmpeg", h.arguments...)

	stdout, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	completeStdout := ""
	lastProgressionTime := h.progressionTime
	for scanner.Scan() {
		dataString := scanner.Text()
		completeStdout += dataString + "\n"
		if duration, err := FindLastDuration(dataString, "Duration: ?"); err == nil {
			h.totalDuration = duration
		}
		if t, err := FindLastDuration(dataString, "out_time=?"); err == nil {
			h.progressionTime = t
		}

		if lastProgressionTime != h.progressionTime {
			handler(h.Progress())
			lastProgressionTime = h.progressionTime
		}
	}

	err = cmd.Wait()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, completeStdout))
	}

	return nil
}

func FindLastDuration(haystack, pattern string) (time.Duration, error) {
	reg, err := regexp.Compile(strings.ReplaceAll(pattern, "?", "(\\d+:\\d+:[\\d\\.]+)"))
	if err != nil {
		return 0, err
	}
	if !reg.MatchString(haystack) {
		return 0, errors.New("not found")
	}

	results := reg.FindAllStringSubmatch(haystack, -1)
	lastMatch := results[len(results)-1][1]
	parts := strings.Split(lastMatch, ":")
	duration, err := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2]))
	if err != nil {
		return 0, err
	}
	return duration, nil
}
