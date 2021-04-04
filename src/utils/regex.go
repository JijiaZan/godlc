package utils

import (
	"regexp"
)

func GetPort(s string) string {
	reg := regexp.MustCompile(`:[0-9]+`)
	res := reg.FindAllStringSubmatch(s, -1)
	return res[0][0][1:]
}

