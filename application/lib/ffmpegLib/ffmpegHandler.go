package ffmpegLib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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
	logrus.Infoln(fmt.Sprintf("Starting %s", reflect.TypeOf(h).String()))
	cmd := exec.Command("ffmpeg", h.arguments...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		stdin, _ := cmd.StdoutPipe()
		io.Copy(ioutil.Discard, stdin)
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	scannerBuffer := make([]byte, 1024*1024*100)
	scanner := bufio.NewScanner(stderr)
	scanner.Buffer(scannerBuffer, len(scannerBuffer)-100)
	scanner.Split(ScanLines)
	completeStdout := ""
	lastProgressionTime := h.progressionTime
	for scanner.Scan() {
		dataString := scanner.Text()
		logrus.Infoln(dataString)
		completeStdout += dataString + "\n"
		if duration, err := FindLastDuration(dataString, "Duration: ?"); err == nil {
			h.totalDuration = duration
		}
		if t, err := FindLastDuration(dataString, "out_time=?"); err == nil {
			h.progressionTime = t
		} else if t, err := FindLastDuration(dataString, "time=?"); err == nil {
			h.progressionTime = t
		}

		if lastProgressionTime != h.progressionTime {
			handler(h.Progress())
			lastProgressionTime = h.progressionTime
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	logrus.Infoln("Waiting for Finish")
	err = cmd.Wait()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, completeStdout))
	}

	return nil
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	nl := bytes.IndexByte(data, '\n')
	cr := bytes.IndexByte(data, '\r')
	if nl >= 0 && (nl < cr && cr >= 0 || cr < 0) {
		// We have a full newline-terminated line.
		return nl + 1, data[0:nl], nil
	}
	if cr >= 0 && (cr <= nl && nl >= 0 || nl < 0) {
		// We have a full newline-terminated line.
		return cr + 1, data[0:cr], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
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
