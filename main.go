// wfs watches the filesystem for changes and automatically runs commands.
package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		printUsage()
		os.Exit(1)
	}

	c, err := ReadConfig(flag.Arg(0))
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
