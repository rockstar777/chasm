package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/stacktic/dropbox"
	"github.com/toqueteos/webbrowser"
	"io/ioutil"
	"path/filepath"
)

type DropboxStore struct {
	Dropbox     dropbox.Dropbox `json:"dropbox"`
	AccessToken string
}

func (d *DropboxStore) Setup() bool {
	// config, err := getConfig()
	db := dropbox.NewDropbox()
	db.SetAppInfo("zpy424sdnluk9c1", "rrmjsz7mlgnholq")

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
	fmt.Println(tok, db.AccessToken())

	// set the oauth info
	d.Dropbox = *db
	d.AccessToken = db.AccessToken()
	fmt.Println("dropbox:", d)

	return true
}

type ClientKey struct {
	Key    string
	Secret string
}

func (d DropboxStore) GetClientKeys() (key, secret string) {
	file, err := ioutil.ReadFile("client/dropbox_secret.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return
	}
	var keys ClientKey
	json.Unmarshal(file, &keys)
	return keys.Key, keys.Secret
}

func (d DropboxStore) Upload(share Share) {
	d.Dropbox.SetAppInfo("zpy424sdnluk9c1", "rrmjsz7mlgnholq")
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
	d.Dropbox.SetAppInfo("zpy424sdnluk9c1", "rrmjsz7mlgnholq")
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

	d.Dropbox.SetAppInfo("zpy424sdnluk9c1", "rrmjsz7mlgnholq")
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

func (d DropboxStore) Restore() string {
	color.Red("not implemented")
	return ""
}

func getDropboxTokenFromWeb() (string, error) {
	authURL := fmt.Sprintf("https://www.dropbox.com/1/oauth2/authorize?client_id=%s&response_type=code", "zpy424sdnluk9c1")
	webbrowser.Open(authURL)

	color.Yellow("Enter Auth Code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		color.Red("Unable to read authorization code %v", err)
		return "", err
	}

	return code, nil
}
