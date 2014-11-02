package imdb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Movie struct {
	Response    string `json:"Response"`
	Error       string `json:"Error"`
	ID          string `json:"imdbID"`
	Title       string `json:"Title"`
	Year        string `json:"Year"`
	Rated       string `json:"Rated"`
	Released    string `json:"Released"`
	Runtime     string `json:"Runtime"`
	Genre       string `json:"Genre"`
	Director    string `json:"Director"`
	Writers     string `json:"Writers"`
	Actors      string `json:"Actors"`
	Description string `json:"Plot"`
	Poster      string `json:"Poster"`
	Rating      string `json:"imdbRating"`
}

func GetInfo(name string, year uint) (*Movie, error) {
	extra := ""
	if year > 0 {
		extra = "&y=" + strconv.FormatUint(uint64(year), 10)
	}
	resp, err := http.Get("http://www.omdbapi.com/?t=" + url.QueryEscape(name) + extra)
	if err != nil {
		return nil, err
	}

	movie := new(Movie)
	err = json.NewDecoder(resp.Body).Decode(movie)
	if err != nil {
		return nil, err
	}

	if movie.Response != "True" {
		return nil, errors.New("Movie not found: " + movie.Error)
	}

	return movie, nil
}
