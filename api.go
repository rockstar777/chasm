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

func loadChasm() {
	CreateOrLoadChasmDir(chasmRoot)
}

func addDropbox(tok string) socketResponse {
	loadChasm()

	var dropbox DropboxStore

	success, message := (&dropbox).Setup(tok)
	if success {
		preferences.DropboxStores = append(preferences.DropboxStores, dropbox)
		preferences.Save()
	}

	return socketResponse{success, message}
}

func addDrive(tok string) socketResponse {
	loadChasm()

	var gdrive GDriveStore

	success, message := (&gdrive).Setup(tok)
	if success {
		preferences.GDriveStores = append(preferences.GDriveStores, gdrive)
		preferences.Save()
	}

	return socketResponse{success, message}
}

func addFolder(path string) socketResponse {
	loadChasm()

	var folderStore FolderStore

	folderStore.Path = path
	success, message := (&folderStore).Setup()

	if success {
		preferences.FolderStores = append(preferences.FolderStores, folderStore)
		preferences.Save()
	}

	return socketResponse{success, message}
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
			sock.Emit("dropbox added", addDropbox(code))
		})

		sock.On("add drive", func(code string) {
			log.Println("got request for add drive", code)
			sock.Emit("drive added", addDrive(code))
		})

		sock.On("add folder", func(path string) {
			log.Println("got request for add folder", path)
			sock.Emit("drive added", addFolder(path))
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
