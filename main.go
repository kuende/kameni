package main

import (
	"runtime"

	"github.com/mailgun/log"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Init([]*log.LogConfig{&log.LogConfig{Name: "console"}})
}
