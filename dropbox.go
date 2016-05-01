package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/stacktic/dropbox"
)

type DropboxStore struct {
	Dropbox     dropbox.Dropbox `json:"dropbox"`
	AccessToken string
	UserID      int
}

func GetClientKeys() (key, secret string) {
	return DropboxClientKey, DropboxClientSecret
}

func (d *DropboxStore) Setup(tok string) (bool, string) {
	db := dropbox.NewDropbox()
	key, secret := GetClientKeys()
	db.SetAppInfo(key, secret)
	db.SetRedirectURL("http://localhost:2000")

	err := db.AuthCode(tok)
	if err != nil {
		color.Red("Unable to get client token: %v", err)
		return false, fmt.Sprintf("Unable to get client token: %v", err)
	}

	account, err := db.GetAccountInfo()
	if err != nil {
		color.Red("Unable to get account information: %v", err)
		return false, fmt.Sprintf("Unable to get account information: %v", err)
	}

	uid := account.UID
	for _, d := range preferences.DropboxStores {
		if d.UserID == uid {
			color.Red("Account for %s already exists.", account.DisplayName)
			return false, fmt.Sprintf("Account for %s already exists.", account.DisplayName)
		}
	}

	// set the oauth info
	d.Dropbox = *db
	d.AccessToken = db.AccessToken()
	d.UserID = uid

	return true, fmt.Sprintf("Success! %v was successfully added.", d.ShortDescription())
}

func (d DropboxStore) Upload(share Share) {
	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)

	fmt.Print(color.MagentaString("Uploading Dropbox/%s...", share.SID))

	input := ioutil.NopCloser(bytes.NewReader(share.Data))
	_, err := d.Dropbox.FilesPut(input, int64(len(share.Data)), string(share.SID), true, "")
	if err != nil {
		color.Red("Error uploading file: ", err)
		return
	}
	//print check mark
	fmt.Print(color.MagentaString("\u2713\n"))
}

func (d DropboxStore) Delete(sid ShareID) {
	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)

	fmt.Print(color.YellowString("Deleting Dropbox/%s...", sid))

	_, err := d.Dropbox.Delete(string(sid))
	if err != nil {
		color.Red("Error deleting file: ", err)
		return
	}

	//print check mark
	fmt.Print(color.GreenString("\u2713\n"))
}

func (d DropboxStore) Description() string {
	label := "Dropbox Store"

	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)

	// get all chasm files from drive
	entry, err := d.Dropbox.Metadata("", true, false, "", "", 0)
	if err != nil {
		color.Red("Unable to iterate names %v", err)
		return label
	}

	account, err := d.Dropbox.GetAccountInfo()
	if err != nil {
		color.Red("Unable to get account information: %v", err)
		return label
	}

	label = fmt.Sprintf("Dropbox Store: %v", account.DisplayName)

	for _, i := range entry.Contents {
		label += fmt.Sprintf("\n\t%s %s", color.YellowString("-"), filepath.Base(i.Path))
	}

	return label
}

func (d DropboxStore) ShortDescription() string {
	label := "Dropbox Store"

	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)

	account, err := d.Dropbox.GetAccountInfo()
	if err != nil {
		color.Red("Unable to get account information: %v", err)
		return label
	}

	return fmt.Sprintf("Dropbox Store: %v", account.DisplayName)
}

func (d DropboxStore) Clean() {
	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)

	entry, err := d.Dropbox.Metadata("", true, false, "", "", 0)
	if err != nil {
		color.Red("Unable to iterate names %v", err)
		return
	}

	for _, i := range entry.Contents {
		name := filepath.Base(i.Path)
		color.Yellow("Removing Dropbox: %v", name)
		d.Dropbox.Delete(name)
	}

	return
}

func (d DropboxStore) Restore() string {
	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)

	restoreDir, err := ioutil.TempDir("", "chasm_dropbox_restore")
	if err != nil {
		color.Red("Error cannot create temp dir: %v", err)
		return ""
	}

	entry, err := d.Dropbox.Metadata("", true, false, "", "", 0)
	if err != nil {
		color.Red("Unable to iterate names %v", err)
		return ""
	}

	color.Yellow("Downloading shares from Dropbox...")

	for _, i := range entry.Contents {
		if !i.IsDir {
			name := filepath.Base(i.Path)
			err := d.Dropbox.DownloadToFile(name, filepath.Join(restoreDir, name), "")
			if err != nil {
				color.Red("Unable to download file %s: %v", name, err)
				return ""
			}
			fmt.Println("\t - got share ", name)
		}
	}

	return restoreDir
}
