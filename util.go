package main

//MARK: Constants
const (
	GoogleDriveClientSecret = string(`{"installed":{"client_id":"713278088797-agohh4u0l5vjscrmn7j0b79i54mtlein.apps.googleusercontent.com","project_id":"peerless-tiger-129119","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://accounts.google.com/o/oauth2/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"wd3-G0rEx4VwuV9gwabi1kbj","redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]}}`)
	DropboxClientKey        = "zpy424sdnluk9c1"
	DropboxClientSecret     = "rrmjsz7mlgnholq"
)

//MARK: Helper Functions

func check(e error) {
	if e != nil {
		panic(e)
	}
}
