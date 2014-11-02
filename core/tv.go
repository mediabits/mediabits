package core

import (
	"bytes"
	"errors"
	"fmt"
	"mediabits/tvdb"
	"strconv"
)

type SeasonInfo struct {
	*sharedInfo
	SeasonNumber int
	Show         *tvdb.Show
	Episodes     []*tvdb.Episode
}

type TVSeason struct {
	Title       string
	Description string
	Info        *SeasonInfo
}

type EpisodeInfo struct {
	*sharedInfo
	Show    *tvdb.Show
	Episode *tvdb.Episode
}

type TVEpisode struct {
	Title       string
	Description string
	Info        *EpisodeInfo
}

func GetShow(name string, year uint) (*tvdb.Show, error) {
	if year > 0 {
		name = name + " (" + strconv.FormatUint(uint64(year), 10) + ")"
	}

	shows, err := tvdb.GetShows(name)
	if err != nil {
		return nil, err
	}
	if len(shows) == 0 {
		return nil, errors.New("No shows found")
	}

	return shows[0], nil
}

func GetSeason(show *tvdb.Show, seasonNumber int, path string, screenshotsCount uint) (*TVSeason, error) {
	episodesChan := make(chan []*tvdb.Episode)
	posterChan := make(chan string)
	errChan := make(chan error)

	// Fetch Episodes
	go (func() {
		episodes, err := show.GetSeasonEpisodes(seasonNumber)
		if err != nil {
			errChan <- err
		} else {
			episodesChan <- episodes
		}
	})()

	// Fetch Poster
	go (func() {
		poster, err := show.GetPoster()
		if err != nil {
			errChan <- err
		} else {
			posterChan <- poster
		}
	})()

	// Get both
	var episodes []*tvdb.Episode
	var poster string
	var err error

	for i := 0; i < 2; i++ {
		select {
		case episodes = <-episodesChan:
			// got episodes
		case poster = <-posterChan:
			// got poster
		case err = <-errChan:
			return nil, err
		}
	}

	// Build the shared info
	sharedInfo, err := getSharedInfo(poster, path, screenshotsCount)
	if err != nil {
		return nil, err
	}

	info := &SeasonInfo{
		sharedInfo,
		seasonNumber,
		show,
		episodes,
	}

	season := &TVSeason{Info: info}
	season.Title = fmt.Sprintf("%s S%02d [ %s / %s / %s / %s / %s ]", info.Show.Name,
		info.SeasonNumber, info.Mediainfo.GeneralSection.Source,
		info.Mediainfo.VideoStream.Format, info.Mediainfo.AudioStream.Format,
		info.Mediainfo.GeneralSection.Container,
		info.Mediainfo.VideoStream.Resolution)

	buf := new(bytes.Buffer)
	RenderText(buf, "tv_season_upload_description", info)
	season.Description = buf.String()

	return season, nil
}

func GetEpisode(show *tvdb.Show, seasonNumber, episodeNumber int, path string, screenshotsCount uint) (*TVEpisode, error) {
	episodeChan := make(chan *tvdb.Episode)
	posterChan := make(chan string)
	errChan := make(chan error)

	// Fetch Episode
	go (func() {
		episode, err := show.GetEpisode(seasonNumber, episodeNumber)
		if err != nil {
			errChan <- err
		} else {
			episodeChan <- episode
		}
	})()

	// Fetch Poster
	go (func() {
		poster, err := show.GetPoster()
		if err != nil {
			errChan <- err
		} else {
			posterChan <- poster
		}
	})()

	// Get both
	var episode *tvdb.Episode
	var poster string
	var err error

	for i := 0; i < 2; i++ {
		select {
		case episode = <-episodeChan:
			// got episode
		case poster = <-posterChan:
			// got poster
		case err = <-errChan:
			return nil, err
		}
	}

	sharedInfo, err := getSharedInfo(poster, path, screenshotsCount)
	if err != nil {
		return nil, err
	}

	info := &EpisodeInfo{
		sharedInfo,
		show,
		episode,
	}

	tvEpisode := &TVEpisode{Info: info}

	tvEpisode.Title = fmt.Sprintf("%s S%02dE%02d [ %s / %s / %s / %s / %s ]",
		info.Show.Name, info.Episode.SeasonNumber, info.Episode.EpisodeNumber,
		info.Mediainfo.GeneralSection.Source, info.Mediainfo.VideoStream.Format,
		info.Mediainfo.AudioStream.Format,
		info.Mediainfo.GeneralSection.Container,
		info.Mediainfo.VideoStream.Resolution)

	buf := new(bytes.Buffer)
	RenderText(buf, "tv_episode_upload_description", info)
	tvEpisode.Description = buf.String()

	return tvEpisode, nil
}
