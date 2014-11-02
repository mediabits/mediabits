package wui

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"mediabits/assets"
	"mediabits/core"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func validate(user, pass string) bool {
	return user == "mediabits" && pass == *password
}

func authReq(f handler) handler {
	if *password != "" {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"mediabits\"")

			if len(r.Header["Authorization"]) == 0 {
				http.Error(w, "authorization failed", http.StatusUnauthorized)
				return
			}

			auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

			if len(auth) != 2 || auth[0] != "Basic" {
				http.Error(w, "bad syntax", http.StatusBadRequest)
				return
			}

			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)

			if len(pair) != 2 || !validate(pair[0], pair[1]) {
				http.Error(w, "authorization failed", http.StatusUnauthorized)
				return
			}

			f(w, r)
		}
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, _ := assets.Asset("html_templates/" + name + ".tpl")
	template.Must(template.New(name).Parse(string(tmpl))).Execute(w, data)
}

func renderStatic(w http.ResponseWriter, name string) {
	asset, err := assets.Asset("html_static/" + name)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("No such file: " + name))
		return
	}

	if strings.HasSuffix(name, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(name, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	}

	w.Write(asset)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		render(w, "index", nil)
	} else if strings.HasPrefix(r.RequestURI, "/static") {
		handleStatic(w, r)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 page not found"))
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.RequestURI, "/", 3)
	if len(parts) != 3 {
		w.WriteHeader(404)
		w.Write([]byte("No such file: " + r.RequestURI))
		return
	}

	renderStatic(w, parts[2])
}

func handleMovie(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid form data (buggy browser/client)"))
		return
	}
	title := r.Form.Get("title")
	yearStr := r.Form.Get("year")
	file := r.Form.Get("file")

	if yearStr == "" {
		yearStr = "0"
	}

	year, err := strconv.ParseUint(yearStr, 10, 32)
	if title == "" || file == "" || err != nil {
		w.WriteHeader(400)
		w.Write([]byte("invalid title/year/file"))
		return
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		w.WriteHeader(400)
		w.Write([]byte("file does not exist"))
		return
	}

	// Get the movie
	movie, err := core.GetMovie(title, uint(year), file, 2)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// Render to response
	core.RenderJSON(w, movie)
}

func handleTV(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid form data (buggy browser/client)"))
		return
	}

	showName := r.Form.Get("show")
	yearStr := r.Form.Get("year")
	seasonNumberStr := r.Form.Get("season")
	episodeNumberStr := r.Form.Get("episode")
	file := r.Form.Get("file")

	if showName == "" || file == "" || seasonNumberStr == "" {
		w.WriteHeader(400)
		w.Write([]byte("show/file cannot be empty"))
		return
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		w.WriteHeader(400)
		w.Write([]byte("file does not exist"))
		return
	}

	var year uint64 = 0
	if yearStr != "" {
		year, err = strconv.ParseUint(yearStr, 10, 64)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("invalid year"))
			return
		}
	}

	show, err := core.GetShow(showName, uint(year))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("failed to fetch show: " + err.Error()))
		return
	}

	seasonNumber, err := strconv.ParseUint(seasonNumberStr, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("invalid season number"))
		return
	}

	if episodeNumberStr == "" { // Season
		season, err := core.GetSeason(show, int(seasonNumber), file, 2)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("failed to fetch season: " + err.Error()))
			return
		}

		core.RenderJSON(w, season)
	} else { // Episode
		episodeNumber, err := strconv.ParseUint(episodeNumberStr, 10, 64)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("invalid episode number"))
			return
		}

		episode, err := core.GetEpisode(show, int(seasonNumber), int(episodeNumber), file, 2)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("failed to fetch episode: " + err.Error()))
			return
		}

		core.RenderJSON(w, episode)
	}
}

type ListFilesResponse struct {
	Directory   string
	Exists      bool
	Directories []string
	Files       []string
}

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		wd, _ := os.Getwd()
		dir = wd
	}
	dir = filepath.Clean(dir)

	res := &ListFilesResponse{
		Directory:   dir,
		Exists:      true,
		Directories: make([]string, 0),
		Files:       make([]string, 0),
	}

	enc := json.NewEncoder(w)

	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		res.Exists = false
		enc.Encode(res)
		return
	}

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			res.Directories = append(res.Directories, entry.Name())
		} else {
			res.Files = append(res.Files, entry.Name())
		}
	}

	enc.Encode(res)
}
