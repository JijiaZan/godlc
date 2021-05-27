package main

import (
	"github.com/JijiaZan/godml/framework"
)

func main() {
	w := &framework.Worker{}
	w.GServe("1888")
}
