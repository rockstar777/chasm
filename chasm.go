package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

/// Chasm Types ///

// CloudStore represents an external cloud storage service that is compatible
// with Chasm
type CloudStore interface {
	Setup()

	Upload(share Share)
	Delete(sid ShareID)

	Description() string
}

// ChasmPref represents user/application preferences
type ChasmPref struct {
	root string

	// the cloud services sharing across
	FolderStores []FolderStore `json:"svcs"`

	// maps files to their shareId
	FileMap map[string]ShareID `json:"files"`
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
	chasmFilePath := p.root + string(filepath.Separator) + chasmPrefFile
	chasmFileBytes, err := json.Marshal(preferences)
	check(err)

	ioutil.WriteFile(chasmFilePath, chasmFileBytes, 0660)
}

/// Chasm Functions ///

var preferences ChasmPref

const chasmPrefFile = ".chasm"
const chasmIgnoreFile = ".chasmignore"

// CreateOrLoadChasmDir creates the root *chasm* folder on the system
// if it does not exist or finds an existing directory
// returns if true if newly created
func CreateOrLoadChasmDir(root string) {
	os.MkdirAll(root, 0777)

	chasmFilePath := path.Join(root, chasmPrefFile)
	chasmFileBytes, err := ioutil.ReadFile(chasmFilePath)
	if err != nil {
		color.Green("Creating new .chasm secure folder")
		preferences.FileMap = make(map[string]ShareID)
		preferences.FileMap[chasmFilePath] = ShareID(chasmPrefFile)
	} else {
		json.Unmarshal(chasmFileBytes, &preferences)
	}

	chasmIgnorePath := path.Join(root, chasmIgnoreFile)
	if _, err := os.Stat(chasmIgnorePath); os.IsNotExist(err) {
		preferences.FileMap[chasmIgnorePath] = ShareID(chasmIgnorePath)
		os.Create(chasmIgnorePath)
	}

	preferences.root = root
	preferences.Save()
}

// IsValidPath checks if a file path is vaild, i.e. it doesn't match any patterns
// in the .chasmignore file
func IsValidPath(filePath string) bool {
	base := filepath.Base(filePath)
	chasmIgnorePath := path.Join(preferences.root, chasmIgnoreFile)
	chasmIgnore, err := os.Open(chasmIgnorePath)
	if err != nil {
		return true
	}

	scanner := bufio.NewScanner(chasmIgnore)
	for scanner.Scan() {
		pattern := scanner.Text()
		// if the file matches anything in .chasmignore, return false
		ok, err := filepath.Match(pattern, base)
		if ok {
			return false
		}
		if err != nil {
			fmt.Println(err)
		}
	}

	// check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return true
}

// AddFile secret shares the file, and uploads each share to corresponding services
// if the file exists already, we delete the remote share first by its shareId
func AddFile(filePath string) {
	if !IsValidPath(filePath) {
		color.Red("Path %s is in .chasmignore. No actions will be performed.", filePath)
		return
	}

	var sid ShareID
	if existingSID, ok := preferences.FileMap[filePath]; ok {
		sid = existingSID
	} else {
		// create unique share_id
		sid = RandomShareID()
		preferences.FileMap[filePath] = sid
	}

	// read the file
	fileBytes, err := ioutil.ReadFile(filePath)
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
func DeleteFile(filePath string) {
	if !IsValidPath(filePath) {
		color.Red("Path %s is in .chasmignore. No actions will be performed.", filePath)
		return
	}

	allCloudStores := preferences.AllCloudStores()

	if sid, ok := preferences.FileMap[filePath]; ok {
		// iteratively delete shares from each cloud store
		for _, cs := range allCloudStores {
			cs.Delete(ShareID(sid))
		}

		delete(preferences.FileMap, filePath)
		preferences.Save()

		color.Green("Deleted share from all cloud stores.")
		return
	}

	color.Red("Path %s is not tracked. Cannot find share id.", filePath)
}
