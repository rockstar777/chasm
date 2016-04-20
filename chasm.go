package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

/// Chasm Types ///

// CloudStore represents an external cloud storage service that is compatible
// with Chasm
type CloudStore interface {
	Upload(share Share)
	Delete(sid ShareID)

	//Restore downloads shares to local restore path
	Restore() string

	Description() string

	Clean()
}

// ChasmPref represents user/application preferences
type ChasmPref struct {
	root string

	// the cloud services sharing across
	FolderStores []FolderStore `json:"folder_stores"`

	// the cloud services sharing across
	GDriveStores []GDriveStore `json:"gdrive_stores"`

	// the cloud services sharing across
	DropboxStores []DropboxStore `json:"dropbox_stores"`

	// maps files to their shareId
	FileMap map[string]ShareID `json:"files"`
}

// RegisteredServices counts all services
func (p ChasmPref) RegisteredServices() int {
	return len(p.FolderStores) + len(p.GDriveStores) + len(p.DropboxStores)
}

// NeedSetup checks if there are enough services to run
func (p ChasmPref) NeedSetup() bool {
	return p.RegisteredServices() < 2
}

// AllCloudStores combines all the cloud stores
func (p ChasmPref) AllCloudStores() []CloudStore {

	// adjust length for new store types
	cloudStores := make([]CloudStore, p.RegisteredServices())

	// all other cloud stores go here
	ind := 0
	for i, fs := range p.FolderStores {
		cloudStores[i] = CloudStore(fs)
		ind += 1
	}

	for j, gds := range p.GDriveStores {
		cloudStores[j+ind] = CloudStore(gds)
		ind += 1
	}

	for k, dbs := range p.DropboxStores {
		cloudStores[k+ind] = CloudStore(dbs)
		ind += 1
	}

	return cloudStores
}

// Save saves the chasm preferences
func (p ChasmPref) Save() {
	chasmFilePath := path.Join(p.root, chasmPrefFile)
	chasmFileBytes, err := json.MarshalIndent(preferences, "", "    ")
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
	_, err = ioutil.ReadFile(chasmIgnorePath)
	if err != nil {
		preferences.FileMap[chasmIgnorePath] = ShareID(chasmIgnoreFile)
		// add *.DS_Store to ignore file by default
		errWrite := ioutil.WriteFile(chasmIgnorePath, []byte(".DS_Store\n"), 0777)
		if errWrite != nil {
			color.Red("Error: could not write to %s: %s", chasmFilePath, errWrite)
		}
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
		color.Blue("Path %s is in .chasmignore. No actions will be performed.", filePath)
		return
	}
	file, _ := os.Open(filePath)
	fi, err := file.Stat()
	if err != nil {
		color.Red("Cannot get file info: %s", err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		files, _ := ioutil.ReadDir(filePath)
		for _, f := range files {
			AddFile(path.Join(filePath, f.Name()))
		}
		return
	case mode.IsRegular():
		break
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

		color.Yellow("Deleted share from all cloud stores.")
		return
	}

	color.Red("Path %s is not tracked. Cannot find share id.", filePath)
}

func DeleteDir(dirPath string) {

	for filePath, _ := range preferences.FileMap {
		dirMatch, _ := path.Split(filePath)
		if path.Clean(dirMatch) != path.Clean(dirPath) {
			continue
		}
		DeleteFile(filePath)
	}
}

// Restore shares to the original files
func Restore() {
	allCloudStores := preferences.AllCloudStores()
	sharePaths := make([]string, len(allCloudStores))

	// (1) first get all shares
	for i, cs := range allCloudStores {
		sp := cs.Restore()
		if sp == "" {
			color.Red("Restore failed for %v", cs)
			return
		}
		sharePaths[i] = sp
	}

	// (2) next restore .chasm file
	chasmFileBytes := restoreShareID(ShareID(chasmPrefFile), sharePaths)

	var restoredPrefs ChasmPref
	err := json.Unmarshal(chasmFileBytes, &restoredPrefs)
	if err != nil {
		color.Red("Cannot restore chasm preferecnes file from cloud services.")
		return
	}

	// (3) finally, for the remaining files, restore and save
	for filePath, sid := range restoredPrefs.FileMap {
		fileBytes := restoreShareID(sid, sharePaths)
		os.MkdirAll(path.Dir(filePath), 0770)
		err := ioutil.WriteFile(filePath, fileBytes, 0770)
		if err != nil {
			color.Red("Error writing restored file %s: %s", filePath, err)
			return
		}
	}
	color.Green("Done. Restored all files!")
}

func restoreShareID(sid ShareID, sharePaths []string) []byte {
	fileShares := make([]Share, len(sharePaths))

	for i, sp := range sharePaths {
		file := path.Join(sp, string(sid))
		dataBytes, err := ioutil.ReadFile(file)
		if err != nil {
			color.Red("(Skipping share) Cannot read file %s: %s", file, err)
			continue
		}

		share := Share{SID: sid, Data: dataBytes}
		fileShares[i] = share
	}

	return CombineShares(fileShares)
}
