package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("error: ")
}

func main() {
	c, err := ReadConfig("config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Watching...")
	watcher := &Watcher{Config: c, Out: os.Stderr}
	watcher.Watch()
}
