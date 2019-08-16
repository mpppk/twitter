package io

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/xerrors"
)

func DownloadEntityMedia(tweet *anaconda.Tweet, entityMedia *anaconda.EntityMedia, index int, downloadDir string) (string, error) {
	mediaRawUrl := entityMedia.Media_url_https
	mediaUrl, err := url.Parse(mediaRawUrl)
	if err != nil {
		return "", xerrors.Errorf("failed to parse media url(%s): %w", mediaRawUrl, err)
	}
	mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
	if len(mediaUrlPaths) == 0 {
		return "", xerrors.Errorf("invalid mediaUrl: %s", mediaRawUrl)
	}
	mediaFileName := mediaUrlPaths[len(mediaUrlPaths)-1]
	fileName := fmt.Sprintf("%d-%d-%s", tweet.Id, index, mediaFileName)
	downloadPath := path.Join(downloadDir, fileName)
	if isExist(downloadPath) {
		return "", nil
	}
	if err := DownloadFile(mediaRawUrl, downloadPath); err != nil {
		return "", xerrors.Errorf("failed to download file to %s", downloadPath)
	}
	return downloadPath, nil
}

func DownloadFile(fileUrl, downloadPath string) (err error) {
	response, err := http.Get(fileUrl)
	if err != nil {
		return xerrors.Errorf("failed to request http get to %s: %w", fileUrl, err)
	}
	defer func() {
		cerr := response.Body.Close()
		if cerr == nil {
			return
		}
		err = xerrors.Errorf("failed to close http response body: %w", err)
	}()

	file, err := os.Create(downloadPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		cerr := file.Close()
		if cerr == nil {
			return
		}
		err = xerrors.Errorf("failed to close file(%s): %w", downloadPath, err)
	}()

	if _, err := io.Copy(file, response.Body); err != nil {
		return xerrors.Errorf("failed to write download file to local file(%s): %w", downloadPath, err)
	}
	return nil
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
