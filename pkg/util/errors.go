package util

import "log"

func LogError(err error) {
	if err != nil {
		log.Printf("LogError: Encountered error %s\n", err)
	}
}
