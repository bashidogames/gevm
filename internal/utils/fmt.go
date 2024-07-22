package utils

import (
	"fmt"
)

func Printlnf(format string, a ...any) (n int, err error) {
	return fmt.Println(fmt.Sprintf(format, a...))
}
