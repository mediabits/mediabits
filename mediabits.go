package main

//go:generate ./assetgen.sh

import (
	"flag"
	"fmt"
	"mediabits/cli"
	"mediabits/ffmpeg"
	"mediabits/mediainfo"
	"mediabits/updater"
	"mediabits/wui"
	"os"
)

var server = flag.Bool("server", false, "run mediabits web user interface")
var checkSystem = flag.Bool("check", false, "check that mediabits and dependencies are working but don't do anything")
var noUpdate = flag.Bool("noupdate", false, "Don't automatically install updates")

func runUpdater(updateConfig *updater.Config, updateChan chan bool) {
	updateInfo, err := updateConfig.CheckForUpdates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check for updates: %s\n", err.Error())
	} else {
		isLatest, err := updateInfo.IsLatest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check latest: %s\n", err.Error())
		} else if !isLatest {
			if *noUpdate {
				fmt.Fprintf(os.Stderr, "Mediabits is out of date please run mediabits without the -noupdate flag to automatically update.\n")
			} else {
				fmt.Fprintf(os.Stderr, "Your mediabits is out of date. Updating...\n")
				fmt.Fprintf(os.Stderr, "DO NOT FORCE QUIT.\n")

				err := updateInfo.Update()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating: %s\n", err.Error())
					updateChan <- false
					return
				}

				fmt.Fprintf(os.Stderr, "The update is complete. Please run mediabits again to use the most recent version.\n")
				updateChan <- true
				return
			}
		}
	}
	updateChan <- false
}

func main() {
	flag.Parse()

	// Auto-update
	updateChan := make(chan bool)
	defer (func() {
		<-updateChan
	})()

	updateConfig := &updater.Config{
		AppName:        "mediabits",
		CurrentVersion: "0.3.3",
		Platforms:      []string{},
		URL:            "https://bbtools.baconseed.org/mediabits/updates",
	}

	go runUpdater(updateConfig, updateChan)

	// Ensure deps are there
	if !ffmpeg.IsInstalled() {
		fmt.Fprintf(os.Stderr, "ffmpeg is not installed\n")
		return
	}

	if !mediainfo.IsInstalled() {
		fmt.Fprintf(os.Stderr, "mediainfo is not installed\n")
		return
	}

	// Run the CLI/WUI respectively
	if *server {
		wui.WuiMain()
	} else if *checkSystem {
		fmt.Println("Success!")
	} else {
		cli.CliMain()
	}
}
