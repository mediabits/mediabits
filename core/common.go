package core

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mediabits/assets"
	"mediabits/ffmpeg"
	"mediabits/images"
	"mediabits/mediainfo"
	"os"
	"strconv"
	"text/template"
)

type sharedInfo struct {
	Image       string
	Screenshots []string
	Mediainfo   *mediainfo.Mediainfo
}

func getSharedInfo(poster, path string, screenshotsCount uint) (*sharedInfo, error) {
	imageURLChan := make(chan string)
	durationChan := make(chan uint)
	fileMediainfoChan := make(chan *mediainfo.Mediainfo)
	errChan := make(chan error)

	// Fetch Image URL
	go (func() {
		imageURL, err := images.UploadImageURL(poster)
		if err != nil {
			errChan <- err
		} else {
			imageURLChan <- imageURL
		}
	})()

	// Get Duration
	go (func() {
		duration, err := ffmpeg.GetDuration(path)
		if err != nil {
			errChan <- err
		} else {
			durationChan <- duration
		}
	})()

	// Get File Mediainfo
	go (func() {
		fileMediainfo, err := mediainfo.GetMediainfo(path)
		if err != nil {
			errChan <- err
		} else {
			fileMediainfoChan <- fileMediainfo
		}
	})()

	// Get all
	var imageURL string
	var duration uint
	var fileMediainfo *mediainfo.Mediainfo
	var err error

	for i := 0; i < 3; i++ {
		select {
		case imageURL = <-imageURLChan:
			// got image url
		case duration = <-durationChan:
			// got duration
		case fileMediainfo = <-fileMediainfoChan:
			// got file mediainfo
		case err = <-errChan:
			return nil, err
		}
	}

	// Take screenshots
	screenshots, err := takeScreenshots(path, duration, screenshotsCount, fileMediainfo)
	if err != nil {
		return nil, err
	}

	return &sharedInfo{
		Image:       imageURL,
		Screenshots: screenshots,
		Mediainfo:   fileMediainfo,
	}, nil
}

func takeScreenshot(path string, point, displayWidth, displayHeight uint) (string, error) {
	// Take screenshot to file
	tmpFile, _ := ioutil.TempFile("", "mbss-"+strconv.FormatUint(uint64(point), 10)+"-")
	name := tmpFile.Name() + ".png"
	tmpFile.Close()
	os.Remove(tmpFile.Name())
	defer os.Remove(name)
	err := ffmpeg.TakeScreenshot(path, name, point, displayWidth, displayHeight)
	if err != nil {
		return "", err
	}

	// Upload screenshot
	url, err := images.UploadImage(name)
	if err != nil {
		return "", err
	}
	return url, nil
}

func takeScreenshots(path string, duration uint, count uint, fileMediainfo *mediainfo.Mediainfo) ([]string, error) {
	point := duration / (count + 1)
	screenshotsChan := make(chan string, count)
	errChan := make(chan error)

	for i := uint(0); i < count; i++ {
		go (func(point uint) {
			screenshot, err := takeScreenshot(path, point, fileMediainfo.VideoStream.DisplayWidth, fileMediainfo.VideoStream.DisplayHeight)
			if err != nil {
				errChan <- err
			} else {
				screenshotsChan <- screenshot
			}
		})(point * (i + 1))
	}

	screenshots := make([]string, 0)
	var err error

	for i := uint(0); i < count; i++ {
		select {
		case screenshot := <-screenshotsChan:
			screenshots = append(screenshots, screenshot)
		case err = <-errChan:
			return nil, err
		}
	}

	return screenshots, nil
}

func Render(out io.Writer, name, tmpl string, data interface{}) {
	template.Must(template.New(name).Parse(tmpl)).Execute(out, data)
}

func RenderText(out io.Writer, name string, data interface{}) {
	tmpl, _ := assets.Asset("templates/" + name + ".txt")
	Render(out, name, string(tmpl), data)
}

func RenderJSON(out io.Writer, data interface{}) {
	json.NewEncoder(out).Encode(data)
}
