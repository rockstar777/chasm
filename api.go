package main

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"path"
	"sync"

	"github.com/googollee/go-socket.io"
)

type socketResponse struct {
	Success bool
	Message string
}

type eventMessage struct {
	Color   string
	Message string
}

// loads chasm information from the preferences file
func loadChasm() {
	CreateOrLoadChasmDir(chasmRoot)
}

// adds a dropbox store
func addDropbox(tok string) socketResponse {
	loadChasm()

	var dropbox DropboxStore

	success, message := (&dropbox).Setup(tok)
	if success {
		preferences.DropboxStores = append(preferences.DropboxStores, dropbox)
		preferences.Save()
	}

	if success {
		messageChannel <- eventMessage{"green", message}
	} else {
		messageChannel <- eventMessage{"red", message}
	}
	return socketResponse{success, message}
}

// adds a google drive store
func addDrive(tok string) socketResponse {
	loadChasm()

	var gdrive GDriveStore

	success, message := (&gdrive).Setup(tok)
	if success {
		preferences.GDriveStores = append(preferences.GDriveStores, gdrive)
		preferences.Save()
	}

	if success {
		messageChannel <- eventMessage{"green", message}
	} else {
		messageChannel <- eventMessage{"red", message}
	}
	return socketResponse{success, message}
}

// adds a folder store
func addFolder(path string) socketResponse {
	loadChasm()

	var folderStore FolderStore

	folderStore.Path = path
	success, message := (&folderStore).Setup()

	if success {
		preferences.FolderStores = append(preferences.FolderStores, folderStore)
		preferences.Save()
	}

	if success {
		messageChannel <- eventMessage{"green", message}
	} else {
		messageChannel <- eventMessage{"red", message}
	}
	return socketResponse{success, message}
}

func cleanChasm() {
	loadChasm()
	if preferences.RegisteredServices() == 0 {
		messageChannel <- eventMessage{"red", "There are no stores to clean."}
		return
	}
	var wg sync.WaitGroup
	for _, cs := range preferences.AllCloudStores() {
		wg.Add(1)
		go func(c CloudStore) {
			defer wg.Done()
			c.Clean()
			messageChannel <- eventMessage{"green", fmt.Sprintf("Done cleaning %v", c.ShortDescription())}
		}(cs)
	}
	wg.Wait()
}

var chasmRoot string
var messageChannel chan eventMessage

func main() {

	usr, _ := user.Current()
	defaultRoot := path.Join(usr.HomeDir, "Chasm")
	chasmRoot = defaultRoot
	messageChannel = make(chan eventMessage)
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(sock socketio.Socket) {
		killChannel := make(chan bool)
		go func(k chan bool) {
			for {
				select {
				case message := <-messageChannel:
					sock.Emit("new event", message)
				case <-k:
					return
				}
			}
		}(killChannel)

		log.Println("on connection")

		// request to get the preferences object
		sock.On("add dropbox", func(code string) {
			log.Println("got request for add dropbox", code)
			sock.Emit("dropbox added", addDropbox(code))
		})

		sock.On("add drive", func(code string) {
			log.Println("got request for add drive", code)
			sock.Emit("drive added", addDrive(code))
		})

		sock.On("add folder", func(path string) {
			log.Println("got request for add folder", path)
			sock.Emit("folder added", addFolder(path))
		})

		sock.On("clean chasm", func() {
			log.Println("got request for clean chasm")
			cleanChasm()
		})

		// on disconnect
		sock.On("disconnection", func() {
			log.Println("on disconnect")
			killChannel <- true
		})
	})
	server.On("error", func(sock socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:4567...")
	log.Fatal(http.ListenAndServe(":4567", nil))
}
