package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"

	"github.com/toqueteos/webbrowser"

	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GDriveStore struct {
	Config     oauth2.Config `json:"oauth_config"`
	OAuthToken oauth2.Token  `json:"oauth_token"`
}

// Setup GDrive
func (g *GDriveStore) Setup() bool {
	config, err := getConfig()

	if err != nil {
		color.Red("Unable to parse client secret file to config: %v", err)
		return false
	}

	tok, err := getTokenFromWeb(config)
	if err != nil {
		color.Red("Unable to get client token: %v", err)
		return false
	}

	// set the oauth info
	g.Config = *config
	g.OAuthToken = *tok

	return true
}

func (g GDriveStore) Upload(share Share) {
	ctx := context.Background()
	config := &g.Config
	client := config.Client(ctx, &g.OAuthToken)

	svc, err := drive.New(client)
	if err != nil {
		color.Red("Unable to retrieve drive Client %v", err)
		return
	}

	// delete existing share
	deleteFilesForShareID(share.SID, svc)

	file := drive.File{}
	now, err := time.Now().MarshalText()
	file.ModifiedTime = string(now)
	file.Name = string(share.SID)
	file.Parents = []string{"appDataFolder"}

	fmt.Println("Uploading: ", file)
	svc.Files.Create(&file).Media(bytes.NewReader(share.Data)).Do()

}
func (g GDriveStore) Delete(sid ShareID) {
	ctx := context.Background()
	config := &g.Config
	client := config.Client(ctx, &g.OAuthToken)

	svc, err := drive.New(client)
	if err != nil {
		color.Red("Unable to retrieve drive Client %v", err)
		return
	}

	// delete existing share
	deleteFilesForShareID(sid, svc)
}

//Restore downloads shares to local restore path
func (g GDriveStore) Restore() string {
	color.Red("Warning: Google Drive Store Unimplemented.")
	return ""
}

func (g GDriveStore) Description() string {
	label := "Google Drive Store"

	ctx := context.Background()
	config := &g.Config
	client := config.Client(ctx, &g.OAuthToken)

	svc, err := drive.New(client)
	if err != nil {
		color.Red("Unable to retrieve drive Client %v", err)
		return label
	}

	// get all chasm files from drive
	r, err := svc.Files.List().Spaces("appDataFolder").Do()
	if err != nil {
		color.Red("Unable to iterate names %v", err)
		return label
	}

	for _, i := range r.Files {
		label += fmt.Sprintf("\n\t%s %s", color.YellowString("-"), i.Name)
	}

	return label
}

/// MARK: Helper Methods ///
func getConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("client/gdrive_client_secret.json")
	if err != nil {
		color.Red("Unable to get client id for chasm %v", err)
		return nil, err
	}

	return google.ConfigFromJSON(b, drive.DriveAppdataScope)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	webbrowser.Open(authURL)

	color.Yellow("Enter Auth Code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		color.Red("Unable to read authorization code %v", err)
		return nil, err
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		color.Red("Unable to retrieve token from web %v", err)
		return tok, err
	}

	return tok, nil
}

func deleteFilesForShareID(sid ShareID, svc *drive.Service) {
	// get all chasm files from drive
	q := fmt.Sprintf("name = '%s'", string(sid))

	r, err := svc.Files.List().Spaces("appDataFolder").Q(q).Do()
	if err != nil {
		color.Red("Unable to search for files to delete: %v", err)
		return
	}

	for _, i := range r.Files {
		svc.Files.Delete(i.Id).Do()
	}
}
