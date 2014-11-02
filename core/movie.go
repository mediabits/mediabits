package core

import (
	"bytes"
	"mediabits/imdb"
)

type MovieInfo struct {
	*sharedInfo
	IMDB *imdb.Movie
}

type Movie struct {
	Title       string
	Description string
	Info        *MovieInfo
}

func GetMovie(name string, year uint, path string, screenshotsCount uint) (*Movie, error) {
	// Get IMDB info
	imdbInfo, err := imdb.GetInfo(name, year)
	if err != nil {
		return nil, err
	}

	sharedInfo, err := getSharedInfo(imdbInfo.Poster, path, screenshotsCount)
	if err != nil {
		return nil, err
	}

	info := &MovieInfo{
		sharedInfo,
		imdbInfo,
	}

	movie := &Movie{Info: info}
	movie.Title = info.IMDB.Title
	buf := new(bytes.Buffer)
	RenderText(buf, "movie_upload_description", info)
	movie.Description = buf.String()

	return movie, nil
}
