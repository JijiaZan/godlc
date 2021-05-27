package utils

import (
	"regexp"
	"strconv"
)

func GetPort(s string) string {
	reg := regexp.MustCompile(`:[0-9]+`)
	res := reg.FindAllStringSubmatch(s, -1)
	return res[0][0][1:]
}

func GetNextPort(s string) string {
	i, _ := strconv.Atoi(GetPort(s))
	return strconv.Itoa(i+1)
}

