package main

import (
	"io"
	"log"
	"path/filepath"

	"golang.org/x/exp/inotify"
)

type Watcher struct {
	Config *Config
	Out    io.Writer
}

func (w *Watcher) Watch() {
	logger := log.New(w.Out, "error: ", 0)

	watcher, err := inotify.NewWatcher()
	if err != nil {
		logger.Panicln(err)
	}

	for _, file := range w.Config.Files {
		err := watcher.AddWatch(
			file.Path,
			inotify.IN_CREATE|inotify.IN_DELETE|inotify.IN_MODIFY,
		)
		if err != nil {
			logger.Println("failed to watch", file.Path)
		}
	}

	for {
		select {
		case ev := <-watcher.Event:
			w.handleEvent(ev)
		case err := <-watcher.Error:
			logger.Println(err)
		}
	}
}

func (w *Watcher) handleEvent(ev *inotify.Event) {
	evDir := filepath.Dir(ev.Name)
	name := filepath.Base(ev.Name)

	for _, file := range w.Config.Files {
		matches, err := filepath.Match(file.Name, name)

		if evDir == file.Path && matches && err == nil {
			log.Println("Event on file", ev.Name)
		}
	}
}
