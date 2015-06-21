package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/inotify"
)

// Watcher is struct that wraps the Linux inotify system.
type Watcher struct {
	Config *Config
	Out    io.Writer
}

// Watch starts watching for the files defined in w.Config and runs the
// appropriate commands when events occur.  This function blocks.
func (w *Watcher) Watch() {
	logger := log.New(w.Out, "", 0)

	watcher, err := inotify.NewWatcher()
	if err != nil {
		logger.Panicln("error:", err)
	}

	for _, file := range w.Config.Files {
		err := watcher.AddWatch(
			file.Path,
			inotify.IN_CREATE|inotify.IN_DELETE|inotify.IN_MODIFY,
		)
		if err != nil {
			logger.Println("error: failed to watch", file.Path)
		}
	}

	for {
		select {
		case ev := <-watcher.Event:
			w.handleEvent(ev)
		case err := <-watcher.Error:
			logger.Println("error:", err)
		}
	}
}

// handleEvent runs the appropriate command based on the event path and type.
func (w *Watcher) handleEvent(ev *inotify.Event) {
	evDir := filepath.Dir(ev.Name)
	name := filepath.Base(ev.Name)

	for _, file := range w.Config.Files {
		// We have to make sure the file we're watching matches the pattern.
		matches, err := filepath.Match(file.Name, name)

		if evDir == file.Path && matches && err == nil {
			env := getEnv(file, ev)

			// Create
			if ev.Mask&inotify.IN_CREATE != 0 && file.Create != "" {
				Run(env, w.Config.Shell, file.Create)
			}

			// Change
			if ev.Mask&inotify.IN_MODIFY != 0 && file.Change != "" {
				Run(env, w.Config.Shell, file.Change)
			}

			// Delete
			if ev.Mask&inotify.IN_DELETE != 0 && file.Delete != "" {
				Run(env, w.Config.Shell, file.Delete)
			}
		}
	}
}

// getEnv returns the environment variables to be passed to the running
// command.  These are the same as the ones for wfs, with the following added:
// - $path: the absolute path to the affected file.
// - $dir: the absolute path to the affected file's parent directory.
// - $dirname: like $dir, but only includes the directory's name, not its path.
// - $filename: the name of the affected file, without its path.
// - $fileradical: like $filename, but without the file extension.
// - $fileext: only the file extension of the affected file (including the dot).
func getEnv(file File, ev *inotify.Event) []string {
	env := os.Environ()
	wd, _ := os.Getwd()

	path := filepath.Join(wd, file.Name)
	dir := filepath.Dir(path)
	dirName := filepath.Base(dir)
	fileName := filepath.Base(ev.Name)

	dotIndex := strings.LastIndex(fileName, ".")
	var fileRadical, fileExt string
	if dotIndex == -1 {
		fileRadical = fileName
		fileExt = ""
	} else {
		fileRadical = fileName[:dotIndex]
		fileExt = fileName[dotIndex:]
	}

	env = append(env, "path="+path)
	env = append(env, "dir="+dir)
	env = append(env, "dirname="+dirName)
	env = append(env, "filename="+fileName)
	env = append(env, "fileradical="+fileRadical)
	env = append(env, "fileext="+fileExt)

	return env
}
