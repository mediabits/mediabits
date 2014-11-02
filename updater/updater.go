package updater

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"bitbucket.org/kardianos/osext"
)

var (
	ErrNoApp      = errors.New("Cannot find update data for that application.")
	ErrNoPlatform = errors.New("Cannot find update data for that platform.")
)

type Artifact struct {
	Version string `json:"version"`
	MD5     string `json:"md5"`
	URL     string `json:"url"`
}

type UpdateInfo struct {
	config    *Config                         `json:"-"`
	Artifacts map[string]map[string]*Artifact `json:"artifacts"`
}

type Config struct {
	AppName        string   `json:"name"`
	CurrentVersion string   `json:"version"`
	Platforms      []string `json:"platforms"`
	URL            string   `json:"url"`
}

func ReadConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	config := new(Config)
	err = json.NewDecoder(configFile).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) CheckForUpdates() (*UpdateInfo, error) {
	resp, err := http.Get(c.URL + "/updater.json")
	if err != nil {
		return nil, err
	}

	updateInfo := new(UpdateInfo)
	err = json.NewDecoder(resp.Body).Decode(updateInfo)
	if err != nil {
		return nil, err
	}

	updateInfo.config = c

	return updateInfo, nil
}

func (u *UpdateInfo) getLatest() (*Artifact, error) {
	platforms, exists := u.Artifacts[u.config.AppName]
	if !exists {
		return nil, ErrNoApp
	}

	artifact, exists := platforms[runtime.GOOS+"_"+runtime.GOARCH]
	if !exists {
		return nil, ErrNoPlatform
	}

	return artifact, nil
}

func (u *UpdateInfo) IsLatest() (bool, error) {
	artifact, err := u.getLatest()
	if err != nil {
		return false, err
	}

	return artifact.Version == u.config.CurrentVersion, nil
}

func (u *UpdateInfo) Update() error {
	exePath, _ := osext.Executable()
	exeDir, _ := osext.ExecutableFolder()
	backupPath := exeDir + string(filepath.Separator) + u.config.AppName + "-backup"
	artifact, err := u.getLatest()
	if err != nil {
		return err
	}

	// Backup the exe

	// Delete the old backup file (if any)
	os.Remove(backupPath)

	err = copyFile(backupPath, exePath)
	if err != nil {
		return err
	}

	// Remove the current exe
	os.Remove(exePath)

	// Download the new exe
	err = artifact.fetch(exePath)
	if err != nil {
		// Attempt to revert to the old version
		os.Remove(exePath)
		copyFile(exePath, backupPath)

		return err
	}

	// Verify the current exe
	isValid, err := artifact.verify(exePath)
	if err != nil {
		// Attempt to revert to the old version
		os.Remove(exePath)
		copyFile(exePath, backupPath)

		return err
	} else if !isValid {
		// Attempt to revert to the old version
		os.Remove(exePath)
		copyFile(exePath, backupPath)

		return errors.New("New binary failed validation.")
	}

	// Fix permissions
	os.Chmod(exePath, 0755)

	return nil
}

func copyFile(to, from string) error {
	// Back up the current binary
	toFile, err := os.Create(to)
	if err != nil {
		return errors.New("Failed to create to file: " + err.Error())
	}
	defer toFile.Close()

	// Copy the current executable to the backup file
	fromFile, err := os.Open(from)
	if err != nil {
		return errors.New("Failed to open from file: " + err.Error())
	}
	defer fromFile.Close()

	// Copy the data
	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return errors.New("Failed to copy file contents: " + err.Error())
	}

	return nil
}

func (a *Artifact) fetch(exePath string) error {
	// Create the binary
	exeFile, err := os.Create(exePath)
	if err != nil {
		return err
	}
	defer exeFile.Close()

	// Fetch the file
	resp, err := http.Get(a.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the file to the binary
	_, err = io.Copy(exeFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func MD5File(path string) (string, error) {
	// Open the binary
	exeFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	// Calculate the sum
	h := md5.New()
	io.Copy(h, exeFile)
	sum := h.Sum(nil)

	return fmt.Sprintf("%x", sum), nil
}

func (a *Artifact) verify(exePath string) (bool, error) {
	sumStr, err := MD5File(exePath)
	if err != nil {
		return false, err
	}

	// Validate
	isValid := sumStr == a.MD5
	return isValid, nil
}
