package utils

import (
	"errors"
	"fmt"
	"os"
)

func DoesExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cannot get path status: %w", err)
	}

	return true, nil
}
