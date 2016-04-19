package main

//MARK: Constants
const (
	GoogleDriveClientSecret = string(`{"installed":{"client_id":"860746801632-2lklu0e3jjv2lqfte5vo2kvlu5cod3g5.apps.googleusercontent.com","project_id":"premium-node-128517","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://accounts.google.com/o/oauth2/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"w9zqAypudT5F2qLA-w6E3oM3","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`)
	DropboxClientKey        = "zpy424sdnluk9c1"
	DropboxClientSecret     = "rrmjsz7mlgnholq"
)

//MARK: Helper Functions

func check(e error) {
	if e != nil {
		panic(e)
	}
}
