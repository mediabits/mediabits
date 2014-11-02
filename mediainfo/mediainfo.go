package mediainfo

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/juju/errgo"
)

var mediainfoBinary = flag.String("mediainfo-bin", "mediainfo", "the path to the mediainfo binary if it is not in the system $PATH")

var completeNameRegexp = regexp.MustCompile("Complete name +: (.*)")

func IsInstalled() bool {
	cmd := exec.Command(*mediainfoBinary)
	err := cmd.Run()
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") ||
			strings.HasSuffix(err.Error(), "executable file not found in %PATH%") ||
			strings.HasSuffix(err.Error(), "executable file not found in $PATH") {
			return false
		} else if strings.HasPrefix(err.Error(), "exit status 255") {
			return true
		}
		fmt.Println("(non-fatal) error determining mediainfo status: " + err.Error())
		fmt.Println("If mediabits does not fail later please report this as a bug.")
	}
	return true
}

func GetMediainfo(mediaFile string) (*Mediainfo, error) {
	cmd := exec.Command(*mediainfoBinary, mediaFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errgo.Notef(err, "failed to open stdout pipe")
	}
	err = cmd.Start()
	if err != nil {
		return nil, errgo.Notef(err, "failed to start")
	}
	scanner := bufio.NewScanner(stdout)
	mediainfo := ""
	for scanner.Scan() {
		line := scanner.Text()
		matches := completeNameRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			line = "Complete name                            : " + filepath.Base(matches[1])
		}
		mediainfo += line + "\n"
	}

	if err = scanner.Err(); err != nil {
		return nil, errgo.Notef(err, "scanner error")
	}

	lexed, err := lex(mediainfo)
	if err != nil {
		return nil, errgo.Notef(err, "error lexing mediainfo")
	}

	parsed, err := parse(lexed)
	if err != nil {
		return nil, errgo.Notef(err, "error parsing mediainfo")
	}

	parsed.Raw = mediainfo

	return parsed, nil
}
