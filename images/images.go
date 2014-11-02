package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	bBImagesBase       = "https://images.baconbits.org/"
	bBImagesUploadPage = bBImagesBase + "upload.php"
	bBImagesImagePath  = bBImagesBase + "images/"
)

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

type BBImagesResponse struct {
	Error    string `json:"error"`
	ErrorMsg string `json:"errorMsg"`
	Image    string `json:"ImgName"`
}

func doImagesRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	imgResp := new(BBImagesResponse)
	err = json.NewDecoder(resp.Body).Decode(imgResp)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	if imgResp.Error == "true" {
		return "", fmt.Errorf("Image Upload error: %s", imgResp.ErrorMsg)
	}

	return bBImagesImagePath + imgResp.Image, nil
}

func UploadImage(filePath string) (string, error) {
	extraParams := map[string]string{}

	req, err := newfileUploadRequest(bBImagesUploadPage, extraParams, "ImageUp", filePath)
	if err != nil {
		return "", err
	}

	return doImagesRequest(req)
}

func UploadImageURL(url string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("url", url)
	err := writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", bBImagesUploadPage, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return doImagesRequest(req)
}
