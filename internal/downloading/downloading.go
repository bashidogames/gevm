package downloading

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bashidogames/gevm/internal/utils"
	"github.com/schollz/progressbar/v3"
)

var ErrNotFound = errors.New("not found")

func Download(url string, path string) error {
	exists, err := utils.DoesExist(path)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		utils.Printlnf("Cached '%s' found", filepath.Base(path))
		return nil
	}

	header, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("failed to request header: %w", err)
	}
	defer header.Body.Close()

	size, err := strconv.ParseInt(header.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse header: %w", err)
	}

	progress := progressbar.NewOptions64(size,
		progressbar.OptionSetDescription(fmt.Sprintf("Downloading '%s'", filepath.Base(path))),
		progressbar.OptionSetWidth(20),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerHead:    "#",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to request download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("download status failure: %w", ErrNotFound)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("download status failure: %s", resp.Status)
	}

	err = os.MkdirAll(filepath.Dir(path), utils.OS_DIRECTORY)
	if err != nil {
		return fmt.Errorf("could not make directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, utils.OS_FILE)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(io.MultiWriter(file, progress), resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	return nil
}

func Fetch(url string, callback func(http.Header, []byte) error) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to request fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("fetch status failure: %w", ErrNotFound)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("fetch status failure: %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = callback(resp.Header, bytes)
	if err != nil {
		return fmt.Errorf("callback failed: %w", err)
	}

	return nil
}
