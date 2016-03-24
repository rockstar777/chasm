package main

import (
	"os"
	"github.com/fatih/color"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"fmt"

)

/// Chasm Types ///

// CloudStore represents an external cloud storage service that is compatible
// with Chasm
type CloudStore interface {
	Setup()

	Upload(share Share)
	Delete(sid ShareID)
}

// ChasmPref represents user/application preferences
type ChasmPref struct {
	root string

	// the cloud services sharing across
	FolderStores []FolderStore	`json:"services"`

	// maps files to their shareId
	FileMap map[string]ShareID	`json:"files"`
}

// RegisteredServices counts all services
func (p ChasmPref) RegisteredServices() int {
	return len(p.FolderStores)
}

// NeedSetup checks if there are enough services to run
func (p ChasmPref) NeedSetup() bool {
	return p.RegisteredServices() < 2
}

// Save saves the chasm preferences
func (p ChasmPref) Save() {
	chasmFilePath := p.root+string(filepath.Separator)+chasmPrefFile
	chasmFileBytes, err := json.Marshal(preferences)
	check(err)

	ioutil.WriteFile(chasmFilePath, chasmFileBytes, 0660)
}


/// Chasm Functions ///

var preferences ChasmPref
const chasmPrefFile = ".chasm"

// CreateOrLoadChasmDir creates the root *chasm* folder on the system
// if it does not exist or finds an existing directory
// returns if true if newly created
func CreateOrLoadChasmDir(root string) {
	os.MkdirAll(root, 0777)

	chasmFilePath := root + chasmPrefFile
	chasmFileBytes, err := ioutil.ReadFile(chasmFilePath)
	if err != nil {
		color.Green("Creating new .chasm secure folder")
		preferences.FileMap = make(map[string]ShareID)
	} else {
		json.Unmarshal(chasmFileBytes, &preferences)
		fmt.Println(preferences.FolderStores[0].Path)
	}

	preferences.root = root
	preferences.Save()
}

// AddFile secret shares the file, and uploads each share to corresponding services
// if the file exists already, we delete the remote share first by its shareId
func AddFile(path string) {

}

// DeleteFile deletes the remote share of this path by its shareId
func DeleteFile(path string) {

}
