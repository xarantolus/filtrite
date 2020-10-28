package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

func DownloadURLs(inputURLs []string, tempDir string) (outputPaths []string, err error) {
	var dlFile = func(url string, file string) (err error) {
		f, err := os.Create(file)
		if err != nil {
			return
		}
		defer f.Close()

		resp, err := httpClient.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}

		_, err = io.Copy(f, resp.Body)

		return
	}

	var errCount int

	for _, dlURL := range inputURLs {
		fn := filepath.Join(tempDir, generateFilename(dlURL))

		err = dlFile(dlURL, fn)
		if err != nil {
			errCount++

			log.Printf("[Warning]: Failed to download %s: %s\n", dlURL, err.Error())
			continue
		}

		outputPaths = append(outputPaths, fn)
	}

	if errCount > (len(inputURLs) / 2) {
		err = fmt.Errorf("%d/%d urls couldn't be downloaded", errCount, len(inputURLs))
	}

	return
}

func generateFilename(url string) string {
	h := sha256.New()
	_, err := h.Write([]byte(url))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil)) + ".txt"
}
