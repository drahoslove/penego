package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

func makeFileWatcher(callback func(string)) (func(string), func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	var currentFile = ""

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if (event.Op & fsnotify.Write) == fsnotify.Write {
					callback(currentFile)
				}
			case err := <-watcher.Errors:
				fmt.Fprintf(os.Stderr, "%s", err)
			}
		}
	}()

	watch := func(file string) {
		if currentFile == file {
			return
		}
		if currentFile != "" {
			err = watcher.Remove(currentFile)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		if file != "" {
			err = watcher.Add(file)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		currentFile = file
	}
	end := func() {
		watcher.Close()
	}

	return watch, end
}
