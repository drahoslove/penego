package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watch  func(string)
	close  func()
	action func()
	isOn   func() bool
}

func makeFileWatcher(callback func(string)) Watcher {
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

	action := func() {
		if currentFile != "" {
			callback(currentFile)
		}
	}
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
		action()
	}
	end := func() {
		watcher.Close()
	}
	isOn := func() bool {
		return currentFile != ""
	}

	return Watcher{watch, end, action, isOn}
}
