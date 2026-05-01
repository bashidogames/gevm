package utils

import (
	"fmt"
)

func Printlnf(format string, a ...any) (n int, err error) {
	return Println(fmt.Sprintf(format, a...))
}

func Println(a ...any) (n int, err error) {
	return fmt.Println(a...)
}
