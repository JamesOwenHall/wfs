package main

import (
	"log"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("error: ")
}

func main() {
	_, err := ReadConfig("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
}
