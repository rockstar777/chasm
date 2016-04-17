package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
)

// FolderStore is a fake cloud store for testing purposes. Simply write
// shares to the folder
type FolderStore struct {
	Path string `json:"path"`
}

// Setup the folder store
func (f FolderStore) Setup() {
	os.MkdirAll(f.Path, 0777)
}

// Upload writes a share to to the folder
func (f FolderStore) Upload(share Share) {
	sharePath := path.Join(f.Path, string(share.SID))
	err := ioutil.WriteFile(sharePath, share.Data, 0770)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	color.Green("Share %s saved successfully!", sharePath)
}

// Delete deletes the share by its shareID
func (f FolderStore) Delete(sid ShareID) {
	sharePath := f.Path + string(sid)
	if _, err := os.Stat(sharePath); err != nil {
		color.Red("Share %s does not exist.", sharePath)
		return
	}

	err := os.Remove(sharePath)
	if err != nil {
		color.Red("Error: could not delete file. %s", err)
		return
	}

	color.Yellow("Share %s deleted successfully!", sid)
}

func (f FolderStore) Description() string {
	return "Folder store at: " + f.Path
}
