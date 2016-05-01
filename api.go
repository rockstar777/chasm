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
		sock.On("add dropbox", func(tok string) {
			log.Println("got request for add dropbox", tok)
			sock.Emit("dropbox added", addDropbox(tok))
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
