package main

import (
	"log"
	"strings"

	"gopkg.in/fsnotify.v1"
)

// StartWatching a path indefinitely.
func StartWatching(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				if strings.Contains(event.Name, ".DS_Store") {
					log.Println("Skipping ignored file.")
					continue
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					AddFile(event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					AddFile(event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					DeleteFile(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	<-done
}
