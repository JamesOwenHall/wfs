package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	c, err := ReadConfig("config.yml")
	if err != nil {
		log.Fatalln("error:", err)
	}

	log.Println("Watching...")
	watcher := &Watcher{Config: c, Out: os.Stderr}
	watcher.Watch()
}
