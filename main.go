// wfs watches the filesystem for changes and automatically runs commands.
package main

import (
	"flag"
	"log"
	"os"
)

const Version = "0.1.0"

func init() {
	log.SetFlags(0)
}

func main() {
	version := flag.Bool("version", false, "Show the version number.")
	flag.Parse()

	if *version == true {
		log.Println("wfs version", Version)
		log.Println("Copyright (c) 2015 James Hall")
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		printUsage()
		os.Exit(1)
	}

	c, err := ReadConfigFile(flag.Arg(0))
	if err != nil {
		log.Fatalln("error:", err)
	}

	log.Println("Watching...")
	watcher := &Watcher{Config: c, Out: os.Stderr}
	watcher.Watch()
}

func printUsage() {
	log.Println("Usage: wfs [filename]")
}
