package cli

import (
	"flag"
	"fmt"
	"mediabits/core"
	"os"
)

// Global
var json = flag.Bool("json", false, "Output as JSON")
var screenshots = flag.Uint("screenshots", 2, "The number of screenshots to be captured")

// Movies
var movieName = flag.String("movie", "", "The name of the movie")
var year = flag.Uint("year", 0, "The year (optional)")

// TV
var showName = flag.String("show", "", "The name of the show")
var seasonNumber = flag.Int("s", -1, "The season")
var episodeNumber = flag.Int("e", -1, "The episode (optional, assumed to be season pack without)")

func CliMain() {
	if flag.Arg(0) == "" {
		fmt.Fprintf(os.Stderr, "You must specify a media file.\n")
		return
	}

	path := flag.Arg(0)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "That file does not exist.\n")
		return
	}

	if *showName != "" { // TV
		show, err := core.GetShow(*showName, *year)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			return
		}

		if *episodeNumber == -1 { // Season
			// Get episode
			season, err := core.GetSeason(show, *seasonNumber, path, *screenshots)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				return
			}

			// Render to stdout
			if *json {
				core.RenderJSON(os.Stdout, season)
			} else {
				core.RenderText(os.Stdout, "tv_season_description", season)
			}
		} else { // Episode
			// Get episode
			episode, err := core.GetEpisode(show, *seasonNumber, *episodeNumber, path, *screenshots)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			// Render to stdout
			if *json {
				core.RenderJSON(os.Stdout, episode)
			} else {
				core.RenderText(os.Stdout, "tv_episode_description", episode)
			}
		}
	} else if *movieName != "" { // Movie
		// Get the movie
		movie, err := core.GetMovie(*movieName, *year, path, *screenshots)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			return
		}

		// Render to stdout
		if *json {
			core.RenderJSON(os.Stdout, movie)
		} else {
			core.RenderText(os.Stdout, "movie_description", movie)
		}

	} else {
		fmt.Fprintf(os.Stderr, "You must specify a movie or a show.\n")
	}
}
