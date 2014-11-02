package mediainfo

import (
	"bufio"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	pixelsRegexp = regexp.MustCompile("(\\d+)")
	dvdripRegexp = regexp.MustCompile("(?i)(DVDRip)")
	hdtvRegexp   = regexp.MustCompile("(?i)(HDTV)")
	webdlRegexp  = regexp.MustCompile("(?i)(WEB)")
	blurayRegexp = regexp.MustCompile("(?i)(Blu-?ray|BDRip)")
)

type General struct {
	FileSize  string
	Container string
	Source    string
	Duration  string
	//BitRate   string
}

type Video struct {
	Format        string
	Resolution    string
	SAR           float64
	DAR           float64
	PAR           float64
	IsAnamorphic  bool
	Width         uint
	Height        uint
	DisplayWidth  uint
	DisplayHeight uint
}

type Audio struct {
	Format string
}

type Mediainfo struct {
	Raw            string
	GeneralSection *General
	VideoStream    *Video
	AudioStream    *Audio
}

func lex(info string) (map[string]map[string]string, error) {
	stringReader := strings.NewReader(info)
	lineScanner := bufio.NewScanner(stringReader)

	data := make(map[string]map[string]string)

	curSection := ""
	sectionData := make(map[string]string)
	for lineScanner.Scan() {
		line := lineScanner.Text()

		// Probably the end of a section
		if strings.TrimSpace(line) == "" {
			if curSection != "" {
				data[curSection] = sectionData
				sectionData = make(map[string]string)
				curSection = ""
			}

			continue
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 1 && curSection == "" { // Section
			curSection = strings.TrimSpace(parts[0])
		} else if len(parts) == 2 { // K/V pair
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			sectionData[key] = value
		} else { // Something is malformed
			continue
		}
	}

	return data, nil
}

func parse(info map[string]map[string]string) (*Mediainfo, error) {
	var err error

	// General
	generalSection := new(General)

	// File Size
	generalSection.FileSize = info["General"]["File size"]

	// Container
	generalSection.Container = info["General"]["Format"]
	if generalSection.Container == "Matroska" {
		generalSection.Container = "MKV"
	} else if generalSection.Container == "MPEG-4" {
		generalSection.Container = "MP4"
	}

	// Source
	generalSection.Source = "Unknown"
	if blurayRegexp.MatchString(info["General"]["Complete name"]) {
		generalSection.Source = "BluRay"
	} else if webdlRegexp.MatchString(info["General"]["Complete name"]) {
		generalSection.Source = "WEB-DL"
	} else if hdtvRegexp.MatchString(info["General"]["Complete name"]) {
		generalSection.Source = "HDTV"
	} else if dvdripRegexp.MatchString(info["General"]["Complete name"]) {
		generalSection.Source = "DVDRip"
	}

	// Duration
	generalSection.Duration = info["General"]["Duration"]

	// BitRate
	// TODO

	// Video
	videoSection := new(Video)

	// Format
	videoSection.Format = info["Video"]["Format"]

	// h.264
	if videoSection.Format == "AVC" {
		writingLibrary, found := info["Video"]["Writing library"]
		if found && strings.HasPrefix(writingLibrary, "x264") {
			videoSection.Format = "x264"
		} else {
			videoSection.Format = "H.264"
		}
	}

	// DivX/XViD
	if videoSection.Format == "MPEG-4 Visual" {
		if info["Video"]["Codec ID"] == "XVID" {
			videoSection.Format = "XVid"
		} else if strings.HasPrefix(info["Video"]["Codec ID/Hint"], "DivX") {
			videoSection.Format = "DivX"
		}
	}

	// Width & Height
	videoSection.Width, _ = dim2pix(info["Video"]["Width"])
	videoSection.Height, _ = dim2pix(info["Video"]["Height"])

	// SAR, DAR and Cthulhu: A love story.
	rawDar := info["Video"]["Display aspect ratio"]
	darArr := strings.SplitN(rawDar, ":", 2)

	if len(darArr) == 2 {
		numerator, err := strconv.ParseFloat(darArr[0], 64)
		if err != nil {
			return nil, errors.New("invalid DAR: num is not a float")
		}

		denominator, err := strconv.ParseFloat(darArr[1], 64)
		if err != nil {
			return nil, errors.New("invalid DAR: den is not a float")
		}

		videoSection.DAR = numerator / denominator
	} else if len(darArr) == 1 {
		videoSection.DAR, err = strconv.ParseFloat(rawDar, 64)
		if err != nil {
			return nil, errors.New("invalid DAR: is not a float: " + darArr[0])
		}
	} else {
		return nil, errors.New("invalid DAR: Display aspect ratio : " + rawDar)
	}

	videoSection.SAR = float64(videoSection.Width) / float64(videoSection.Height)
	videoSection.PAR = videoSection.DAR / videoSection.SAR

	if math.Abs(videoSection.SAR-videoSection.DAR) < 0.1 { // Not Anamorphic
		videoSection.IsAnamorphic = false
		videoSection.DisplayWidth = videoSection.Width
		videoSection.DisplayHeight = videoSection.Height
	} else if videoSection.SAR < videoSection.DAR { // Horizontal Anamorphic
		videoSection.IsAnamorphic = true
		videoSection.DisplayWidth = uint(float64(videoSection.Height) * videoSection.DAR)
		videoSection.DisplayHeight = videoSection.Height
	} else { // Vertical anamorphic
		videoSection.IsAnamorphic = true
		videoSection.DisplayWidth = videoSection.Width
		videoSection.DisplayHeight = uint(float64(videoSection.Width) * videoSection.DAR)
	}

	// Resolution
	videoSection.Resolution = "Unknown"
	if videoSection.DisplayWidth < 640 && videoSection.DisplayHeight < 480 {
		videoSection.Resolution = "SD"
	} else if videoSection.DisplayWidth <= 640 && videoSection.DisplayHeight <= 480 {
		videoSection.Resolution = "480p"
	} else if videoSection.DisplayWidth <= 1024 && videoSection.DisplayHeight <= 576 {
		videoSection.Resolution = "576p"
	} else if videoSection.DisplayWidth <= 1280 && videoSection.DisplayHeight <= 720 {
		videoSection.Resolution = "720p"
	} else if videoSection.DisplayWidth <= 1920 && videoSection.DisplayHeight <= 1080 {
		videoSection.Resolution = "1080p"
	}

	// Audio
	audioSectionName := ""
	if _, found := info["Audio"]; found {
		audioSectionName = "Audio"
	} else if _, found := info["Audio #1"]; found {
		audioSectionName = "Audio #1"
	} else {
		return nil, errors.New("no audio section found")
	}

	audioSection := new(Audio)

	// Format
	audioSection.Format = info[audioSectionName]["Format"]

	// MP3
	if audioSection.Format == "MPEG Audio" && info[audioSectionName]["Format profile"] == "Layer 3" {
		audioSection.Format = "MP3"
	}

	return &Mediainfo{
		Raw:            "",
		GeneralSection: generalSection,
		VideoStream:    videoSection,
		AudioStream:    audioSection,
	}, nil
}

func dim2pix(in string) (uint, error) {
	in = strings.TrimSpace(in)
	in = strings.Replace(in, " ", "", 1)
	if str := pixelsRegexp.FindString(in); len(str) > 0 {
		n, err := strconv.ParseUint(str, 10, 32)
		return uint(n), err
	}
	n, err := strconv.ParseUint(in, 10, 32)
	return uint(n), err
}
