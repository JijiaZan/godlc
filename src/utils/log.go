package utils

import (
	"log"
)

const DEBUG = true

func DPrintf(format string, v ...interface{}) {
	if DEBUG {
		log.Printf(format+"\n", v...)
	}
}