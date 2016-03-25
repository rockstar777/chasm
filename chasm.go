package main

import (
	"os"
	"github.com/fatih/color"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
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

// AllCloudStores combines all the cloud stores
func (p ChasmPref) AllCloudStores() []CloudStore {

	// adjust length for new store types
	cloudStores := make([]CloudStore, len(p.FolderStores))

	// all other cloud stores go here
	for i, fs := range p.FolderStores {
		cloudStores[i] = CloudStore(fs)
	}

	return cloudStores
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
		preferences.FileMap[chasmFilePath] = ShareID(chasmPrefFile)
	} else {
		json.Unmarshal(chasmFileBytes, &preferences)
	}

	preferences.root = root
	preferences.Save()
}

// AddFile secret shares the file, and uploads each share to corresponding services
// if the file exists already, we delete the remote share first by its shareId
func AddFile(path string) {

	var sid ShareID
	if path == preferences.root + chasmPrefFile {
		// if path is the .chasm, use the const sid
		sid = ShareID(".chasm")
	} else {
		// create unique share_id
		sid = RandomShareID()
		preferences.FileMap[path] = sid
	}


	// read the file
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red("Cannot read file: %s", err)
		return
	}

	// create the shares
	allCloudStores := preferences.AllCloudStores()
	shares := CreateShares(fileBytes, sid, len(allCloudStores))

	// iteratively upload shares with each cloud store
	for i, cs := range allCloudStores {
		cs.Upload(shares[i])
	}

	// only save pref if it's not a .chasm
	if sid != ShareID(".chasm") {
		preferences.Save()
	}

}

// DeleteFile deletes the remote share of this path by its shareId
func DeleteFile(path string) {
	allCloudStores := preferences.AllCloudStores()

	if sid, ok := preferences.FileMap[path]; ok {
		// iteratively delete shares from each cloud store
		for _, cs := range allCloudStores {
			cs.Delete(ShareID(sid))
		}

		delete(preferences.FileMap, path)
		preferences.Save()

		color.Green("Deleted share from all cloud stores.")
		return
	}

	color.Red("Path %s is not tracked. Cannot find share id.", path)
}
