package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/stacktic/dropbox"
	"github.com/toqueteos/webbrowser"
)

type DropboxStore struct {
	Dropbox     dropbox.Dropbox `json:"dropbox"`
	AccessToken string
}

type ClientKey struct {
	Key    string
	Secret string
}

func GetClientKeys() (key, secret string) {
	file, err := ioutil.ReadFile("client/dropbox_secret.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return
	}
	var keys ClientKey
	json.Unmarshal(file, &keys)
	return keys.Key, keys.Secret
}

func (d *DropboxStore) Setup() bool {
	// config, err := getConfig()
	db := dropbox.NewDropbox()
	key, secret := GetClientKeys()
	db.SetAppInfo(key, secret)

	tok, err := getDropboxTokenFromWeb()
	if err != nil {
		color.Red("Unable to get client token: %v", err)
		return false
	}

	err = db.AuthCode(tok)
	if err != nil {
		color.Red("Unable to get client token: %v", err)
		return false
	}

	// set the oauth info
	d.Dropbox = *db
	d.AccessToken = db.AccessToken()

	return true
}

func (d DropboxStore) Upload(share Share) {
	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)
	fmt.Printf("Uploading %s to Dropbox...\n", share.SID)
	input := ioutil.NopCloser(bytes.NewReader(share.Data))
	_, err := d.Dropbox.FilesPut(input, int64(len(share.Data)), string(share.SID), true, "")
	if err != nil {
		color.Red("Error uploading file: ", err)
		return
	}
	fmt.Printf("Uploaded %s to Dropbox!\n", share.SID)
}

func (d DropboxStore) Delete(sid ShareID) {
	key, secret := GetClientKeys()
	d.Dropbox.SetAppInfo(key, secret)
	d.Dropbox.SetAccessToken(d.AccessToken)
	fmt.Printf("Deleting %s from Dropbox...\n", sid)
	_, err := d.Dropbox.Delete(string(sid))
	if err != nil {
		color.Red("Error deleting file: ", err)
		return
	}
	fmt.Printf("Deleted %s from Dropbox!\n", sid)
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

	for _, i := range entry.Contents {
		label += fmt.Sprintf("\n\t%s %s", color.YellowString("-"), filepath.Base(i.Path))
	}

	return label
}

func (d DropboxStore) Clean() {
	color.Yellow("Cleaning dropbox:")

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
		fmt.Println("\t- remove ", name)
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

func getDropboxTokenFromWeb() (string, error) {
	key, _ := GetClientKeys()
	authURL := fmt.Sprintf("https://www.dropbox.com/1/oauth2/authorize?client_id=%s&response_type=code", key)
	webbrowser.Open(authURL)

	color.Yellow("Enter Auth Code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		color.Red("Unable to read authorization code %v", err)
		return "", err
	}

	return code, nil
}
