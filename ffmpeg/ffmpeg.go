package ffmpeg

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/juju/errgo"
)

var ffmpegBinary = flag.String("ffmpeg-bin", "ffmpeg", "the path to the mediainfo binary if it is not in the system $PATH")

var durationRegexp = regexp.MustCompile(".*Duration: (\\d{2}):(\\d{2}):(\\d{2}).*")

func IsInstalled() bool {
	cmd := exec.Command(*ffmpegBinary)
	err := cmd.Run()
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") ||
			strings.HasSuffix(err.Error(), "executable file not found in %PATH%") ||
			strings.HasSuffix(err.Error(), "executable file not found in $PATH") {
			return false
		} else if strings.HasPrefix(err.Error(), "exit status 1") {
			return true
		}
		fmt.Println("(non-fatal) error determining ffmpeg status: " + err.Error())
		fmt.Println("If mediabits does not fail later please report this as a bug.")
	}
	return true
}

func GetDuration(mediaFile string) (uint, error) {
	cmd := exec.Command(*ffmpegBinary, "-i", mediaFile)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, errgo.Notef(err, "failed to open stderr pipe")
	}
	err = cmd.Start()
	if err != nil {
		return 0, errgo.Notef(err, "failed to start")
	}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		matches := durationRegexp.FindStringSubmatch(line)
		if len(matches) == 4 {
			var hours uint64
			var minutes uint64
			var seconds uint64
			hours, err = strconv.ParseUint(matches[1], 10, 32)
			if err == nil {
				minutes, err = strconv.ParseUint(matches[2], 10, 32)
			}
			if err == nil {
				seconds, err = strconv.ParseUint(matches[3], 10, 32)
			}
			if err == nil {
				return uint(hours*60*60 + minutes*60 + seconds), nil
			}
		}
	}

	return 0, scanner.Err()
}

func TakeScreenshot(mediaFile, outFile string, at uint, width uint, height uint) error {
	cmd := exec.Command(*ffmpegBinary,
		"-ss", strconv.FormatUint(uint64(at), 10),
		"-i", mediaFile,
		"-f", "image2",
		"-vframes", "1",
		"-vf", "scale="+strconv.FormatUint(uint64(width), 10)+":"+strconv.FormatUint(uint64(height), 10),
		outFile)
	err := cmd.Run()
	return err
}
