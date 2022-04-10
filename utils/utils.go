package utils

import (
	"fmt"
	"os"
)

func CheckError(err error, exitWhenErr bool) bool {
	if err != nil {
		fmt.Println("Err: ", err)
		if exitWhenErr {
			os.Exit(1)
		}
		return true
	}
	return false
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
