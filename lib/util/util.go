package util

import (
	"fmt"
	"os"
)

func IsExistDir(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		return true
	}

	return false
}

func GetURL(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
