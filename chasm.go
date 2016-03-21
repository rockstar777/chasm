package main

/// Chasm Types ///

// ShareID is a uniqiue id to represent uploaded shares
type ShareID string

// Share represents a secret share of a file
type Share struct {
	SID  ShareID
	Data []byte
}

// CloudStore represents an external cloud storage service that is compatible
// with Chasm
type CloudStore interface {
	Setup(username, password string)
	SetupInteractively()

	Upload(share Share)
	Delete(sid ShareID)
}

// ChasmPref represents user/application preferences
type ChasmPref struct {
	// the cloud services sharing across
	RegisteredServices []CloudStore

	// the chasm secure directory where everything happens
	ChasmDir string

	// maps files to their shareId
	FileMap map[string]ShareID
}

/// Chasm Functions ///

var preferences ChasmPref

// CreateOrLoadChasmDir creates the root *chasm* folder on the system
// if it does not exist or finds an existing directory
// returns if true if newly created
func CreateOrLoadChasmDir(root string) bool {
	preferences = ChasmPref{}

	return true
}

// RunSetup runs the setup wizard to get started
func RunSetup() {

}

// AddFile secret shares the file, and uploads each share to corresponding services
// if the file exists already, we delete the remote share first by its shareId
func AddFile(path string) {

}

// DeleteFile deletes the remote share of this path by its shareId
func DeleteFile(path string) {

}
