package main

import (
	"log"

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
				if event.Op&fsnotify.Write == fsnotify.Write {
					AddFile(event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					DeleteFile(event.Name)
				}
				//TODO: Handle rename, create, etc...
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	<-done
}
