package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"sync"

	"github.com/fatih/color"
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
			color.Green("Done cleaning %v", c.ShortDescription())
			messageChannel <- eventMessage{"green", fmt.Sprintf("Done cleaning %v", c.ShortDescription())}
		}(cs)
	}
	wg.Wait()
}

func syncChasm() {
	messageChannel <- eventMessage{"yellow", "Cleaning chasm..."}
	cleanChasm()
	color.Green("Done cleaning.\nBeginning sync:")
	messageChannel <- eventMessage{"green", "Done cleaning. Beginning Sync!"}

	if preferences.NeedSetup() {
		color.Red("Error: not enough services. Cannot sync.")
		messageChannel <- eventMessage{"green", "Error: not enough services. Cannot sync."}
		return
	}

	files, _ := ioutil.ReadDir(preferences.root)
	currentFileMap := make(map[string]bool)
	for _, f := range files {
		if f.Name() == chasmPrefFile {
			continue
		}
		path := path.Join(preferences.root, f.Name())
		currentFileMap[path] = true
		fmt.Println("Sharing ", path)
		messageChannel <- eventMessage{"black", fmt.Sprintf("Sharing %v...", path)}
		AddFile(path)
	}

	// remove invalid entries in existing file map
	for filePath, _ := range preferences.FileMap {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			delete(preferences.FileMap, filePath)
		}
	}

	preferences.Save()
	AddFile(path.Join(preferences.root, chasmPrefFile))

	color.Green("Done syncing.")
	messageChannel <- eventMessage{"green", "Done syncing!"}
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
			log.Println("Starting to listen on messageChannel")
			for {
				select {
				case message := <-messageChannel:
					sock.Emit("new event", message)
				case <-k:
					log.Println("Killing listener on messageChannel")
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
			sock.Emit("chasm cleaned")
		})

		sock.On("sync chasm", func() {
			log.Println("got request for clean chasm")
			syncChasm()
			sock.Emit("chasm synced")
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
