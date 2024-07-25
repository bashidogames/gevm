package archiving

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bashidogames/gdvm/internal/utils"
)

func Unzip(from string, to string) error {
	reader, err := zip.OpenReader(from)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer reader.Close()

	utils.Printlnf("Unzipping '%s'", filepath.Base(from))

	err = unzip(reader, to)
	if err != nil {
		return fmt.Errorf("cannot unzip file: %w", err)
	}

	return nil
}

func unzip(reader *zip.ReadCloser, to string) error {
	for _, file := range reader.File {
		path := filepath.Join(to, file.Name)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, utils.OS_DIRECTORY)
			if err != nil {
				return fmt.Errorf("could not make directory: %w", err)
			}

			continue
		}

		err := os.MkdirAll(filepath.Dir(path), utils.OS_DIRECTORY)
		if err != nil {
			return fmt.Errorf("could not make directory: %w", err)
		}

		dst, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode())
		if err != nil {
			return fmt.Errorf("could not create destination file: %w", err)
		}
		defer dst.Close()

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("could not open zip file: %w", err)
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return fmt.Errorf("could not copy file: %w", err)
		}
	}

	return nil
}
