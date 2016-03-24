package main

// FolderStore is a fake cloud store for testing purposes. Simply write
// shares to the folder
type FolderStore struct {
    Path string `json:"path"`
}

// Setup the folder store
func (f FolderStore) Setup() {

}

// Upload writes a share to to the folder
func (f FolderStore) Upload(share Share) {

}

// Delete deletes the share by its shareID
func (f FolderStore) Delete(sid ShareID) {

}
