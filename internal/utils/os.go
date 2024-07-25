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

const OS_DIRECTORY = OS_USER_RWX | OS_GROUP_RX | OS_OTHER_RX | os.ModeDir
const OS_EXECUTABLE = OS_USER_RWX | OS_GROUP_RX | OS_OTHER_RX
const OS_FILE = OS_USER_RW | OS_GROUP_R | OS_OTHER_R

const OS_USER_RWX = 0700
const OS_USER_RW = 0600
const OS_USER_RX = 0500
const OS_USER_R = 0400
const OS_USER_W = 0200
const OS_USER_X = 0100

const OS_GROUP_RWX = 0070
const OS_GROUP_RW = 0060
const OS_GROUP_RX = 0050
const OS_GROUP_R = 0040
const OS_GROUP_W = 0020
const OS_GROUP_X = 0010

const OS_OTHER_RWX = 0007
const OS_OTHER_RW = 0006
const OS_OTHER_RX = 0005
const OS_OTHER_R = 0004
const OS_OTHER_W = 0002
const OS_OTHER_X = 0001
