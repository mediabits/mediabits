package tvdb

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
)

const (
	baseURL    = "http://thetvdb.com/api/"
	bannersURL = "http://thetvdb.com/banners/"
	apiKey     = "0629B785CE550C8D"
	apiURL     = baseURL + apiKey + "/"
)

type Show struct {
	XMLName     xml.Name `xml:"Series"`
	ID          string   `xml:"id"`
	Name        string   `xml:"SeriesName"`
	Description string   `xml:"Overview"`
	FirstAired  string   `xml:"FirstAired"`
	Network     string   `xml:"Network"`
	IMDBID      string   `xml:"IMDB_ID"`
}

type ShowResponse struct {
	XMLName xml.Name `xml:"Data"`
	Shows   []*Show  `xml:"Series"`
}

type Episode struct {
	XMLName       xml.Name `xml:"Episode"`
	ID            string   `xml:"id"`
	Name          string   `xml:"EpisodeName"`
	Description   string   `xml:"Overview"`
	SeasonNumber  int      `xml:"SeasonNumber"`
	EpisodeNumber int      `xml:"EpisodeNumber"`
	FirstAired    string   `xml:"FirstAired"`
	Rating        string   `xml:"Rating"`
	Director      string   `xml:"Director"`
	Writers       string   `xml:"Writer"`
}

type EpisodesResponse struct {
	XMLName  xml.Name   `xml:"Data"`
	Episodes []*Episode `xml:"Episode"`
}

type Banner struct {
	XMLName xml.Name `xml:"Banner"`
	Type    string   `xml:"BannerType"`
	Path    string   `xml:"BannerPath"`
}

type BannersResponse struct {
	XMLName xml.Name  `xml:"Banners"`
	Banners []*Banner `xml:"Banner"`
}

func GetShows(name string) ([]*Show, error) {
	resp, err := http.Get(baseURL + "GetSeries.php?seriesname=" + url.QueryEscape(name))
	if err != nil {
		return nil, err
	}

	showResp := new(ShowResponse)
	err = xml.NewDecoder(resp.Body).Decode(showResp)
	if err != nil {
		return nil, err
	}

	return showResp.Shows, nil
}

func (s *Show) GetEpisodes() ([]*Episode, error) {
	resp, err := http.Get(apiURL + "series/" + s.ID + "/all/en.xml")
	if err != nil {
		return nil, err
	}

	epsResp := new(EpisodesResponse)
	err = xml.NewDecoder(resp.Body).Decode(epsResp)
	if err != nil {
		return nil, err
	}

	return epsResp.Episodes, nil
}

func (s *Show) GetSeasonEpisodes(seasonNumber int) ([]*Episode, error) {
	episodes, err := s.GetEpisodes()
	if err != nil {
		return nil, err
	}
	seasonEpisodes := make([]*Episode, 0)

	for _, episode := range episodes {
		if episode.SeasonNumber == seasonNumber {
			seasonEpisodes = append(seasonEpisodes, episode)
		}
	}

	if len(seasonEpisodes) == 0 {
		return nil, errors.New("no episodes found")
	} else {
		return seasonEpisodes, nil
	}
}

func (s *Show) GetEpisode(seasonNumber, episodeNumber int) (*Episode, error) {
	episodes, err := s.GetEpisodes()
	if err != nil {
		return nil, err
	}

	for _, episode := range episodes {
		if episode.SeasonNumber == seasonNumber && episode.EpisodeNumber == episodeNumber {
			return episode, nil
		}
	}

	return nil, errors.New("episode not found")
}

func (s *Show) GetPoster() (string, error) {
	resp, err := http.Get(apiURL + "series/" + s.ID + "/banners.xml")
	if err != nil {
		return "", err
	}

	bannersResp := new(BannersResponse)
	err = xml.NewDecoder(resp.Body).Decode(bannersResp)
	if err != nil {
		return "", err
	}

	for _, banner := range bannersResp.Banners {
		if banner.Type == "poster" {
			return bannersURL + banner.Path, nil
		}
	}

	return "https://upload.wikimedia.org/wikipedia/en/f/f9/No-image-available.jpg", nil
}
