package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/users"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/fatih/color"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type DropboxStore struct {
	Dropbox     users.Client
	AccessToken string
	UserID      string
	config dropbox.Config
}



func GetClientKeys() (key, secret string) {
	return DropboxClientKey, DropboxClientSecret
}

func connect(key,secret,code string)  *oauth2.Token{
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     key,
		ClientSecret: secret,
		Endpoint:  dropbox.OAuthEndpoint(""),
			}
	tok, err := conf.Exchange(ctx, code)
	return tok 
}

func (d *DropboxStore) Setup() bool {
	
	key, secret := GetClientKeys()
	
	code, err := getDropboxTokenFromWeb()
	tok := connect(key,secret,code)
	config := dropbox.Config{
      		Token: tok.AccessToken,
      		LogLevel: dropbox.LogOff,
  				}
	db := users.New(config)
	
	if err != nil {
		color.Red("Unable to get client token: %v", err)
		return false
	}

	

	account, err := db.GetCurrentAccount()
	if err != nil {
		color.Red("Unable to get account information: %v", err)
		return false
	}

	uid := account.AccountId
	for _, d := range preferences.DropboxStores {
		if d.UserID == uid {
			color.Red("Account for %s already exists.", account.Name.DisplayName)
			return false
		}
	}

	
	d.Dropbox = db
	d.AccessToken = tok.AccessToken
	d.UserID = uid
	d.config = config

	return true
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
