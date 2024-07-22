package utils

import (
	"fmt"
	"strconv"
)

func AtoiIfNotEmpty(value string) (int, error) {
	if len(value) == 0 {
		return 0, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("conversion failed: %w", err)
	}

	return result, nil
}
