package main

import (
	"crypto/sha256"
	"encoding/base64"
)

//MARK: Constants
const (
	GoogleDriveClientSecret = string(`{"installed":{"client_id":"713278088797-agohh4u0l5vjscrmn7j0b79i54mtlein.apps.googleusercontent.com","project_id":"peerless-tiger-129119","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://accounts.google.com/o/oauth2/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"wd3-G0rEx4VwuV9gwabi1kbj","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`)
	DropboxClientKey        = "0jz22lrmv7v0tkw"
	DropboxClientSecret     = "axs0n3htxsn6o5f"
)

//MARK: Helper Functions

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// MARK: SHA256 Helpers

func SHA256Base64URL(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func checkSHA2(hash string, data []byte) bool {
	return SHA256Base64URL(data) == hash
}
