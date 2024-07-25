package utils

import "fmt"

func SelectFirstNotEmpty(values ...string) (string, error) {
	for _, value := range values {
		if len(value) > 0 {
			return value, nil
		}
	}

	return "", fmt.Errorf("all values empty")
}
