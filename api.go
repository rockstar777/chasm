package main

import (
	"log"
	"net/http"
	"os/user"
	"path"

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

func loadChasm() {
	CreateOrLoadChasmDir(chasmRoot)
}

func addDropbox(tok string) (socketResponse, eventMessage) {
	loadChasm()

	var dropbox DropboxStore

	success, message := (&dropbox).Setup(tok)
	if success {
		preferences.DropboxStores = append(preferences.DropboxStores, dropbox)
		preferences.Save()
	}

	if success {
		return socketResponse{success, message}, eventMessage{"green", message}
	}
	return socketResponse{success, message}, eventMessage{"red", message}
}

func addDrive(tok string) (socketResponse, eventMessage) {
	loadChasm()

	var gdrive GDriveStore

	success, message := (&gdrive).Setup(tok)
	if success {
		preferences.GDriveStores = append(preferences.GDriveStores, gdrive)
		preferences.Save()
	}

	if success {
		return socketResponse{success, message}, eventMessage{"green", message}
	}
	return socketResponse{success, message}, eventMessage{"red", message}
}

func addFolder(path string) (socketResponse, eventMessage) {
	loadChasm()

	var folderStore FolderStore

	folderStore.Path = path
	success, message := (&folderStore).Setup()

	if success {
		preferences.FolderStores = append(preferences.FolderStores, folderStore)
		preferences.Save()
	}

	if success {
		return socketResponse{success, message}, eventMessage{"green", message}
	}
	return socketResponse{success, message}, eventMessage{"red", message}
}

var chasmRoot string

func main() {

	usr, _ := user.Current()
	defaultRoot := path.Join(usr.HomeDir, "Chasm")
	chasmRoot = defaultRoot

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(sock socketio.Socket) {
		log.Println("on connection")

		// request to get the preferences object
		sock.On("add dropbox", func(code string) {
			log.Println("got request for add dropbox", code)
			response, message := addDropbox(code)
			sock.Emit("dropbox added", response)
			sock.Emit("new event", message)
		})

		sock.On("add drive", func(code string) {
			log.Println("got request for add drive", code)
			response, message := addDrive(code)
			sock.Emit("dropbox added", response)
			sock.Emit("new event", message)
		})

		sock.On("add folder", func(path string) {
			log.Println("got request for add folder", path)
			response, message := addFolder(path)
			sock.Emit("dropbox added", response)
			sock.Emit("new event", message)
		})

		// on disconnect
		sock.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(sock socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
