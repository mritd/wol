package utils

import (
	"log"
	"os"
)

func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CheckAndExit(err error) {
	if !CheckErr(err) {
		os.Exit(1)
	}
}
