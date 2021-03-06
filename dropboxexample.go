package main


import (
	"fmt"
	"errors"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/users"	
	"os"
	"path/filepath"
	"io/ioutil"
	"text/tabwriter"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"github.com/mitchellh/ioprogress"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"path"
	"io"
	"github.com/dustin/go-humanize"
	"time"
	"github.com/toqueteos/webbrowser"
	"strings"
	
)



var config dropbox.Config
type TokenMap map[string]map[string]string
const chunkSize int64 = 1 << 24

func check(e error){
	if e!=nil {
		panic(e)
}
}

func readtoken(filePath string)(string ,error){
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}



	return string(b), nil
	
}


func writetoken(filePath string,p string){
		

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Doesn't exist; lets create it
		err = os.MkdirAll(filepath.Dir(filePath), 0700)
		check(err)
	}
	b:=[]byte(p)
	err := ioutil.WriteFile(filePath, b, 0600)
	check(err)
	
}
func confffiggg() {
filePath := path.Join("/home/rockstar", "dbxcli", "test.txt")

ctx := context.Background()
conf := &oauth2.Config{
    ClientID:     "0jz22lrmv7v0tkw",
    ClientSecret: "axs0n3htxsn6o5f",
    Endpoint: dropbox.OAuthEndpoint(""),
}
token,err := readtoken(filePath);if err!=nil{
panic(err)
}
if token=="" {

url := conf.AuthCodeURL("state")
fmt.Printf("Visit the URL for the auth dialog: %v", url)
webbrowser.Open(url)


var code string
if _, err := fmt.Scan(&code); err != nil {
    log.Fatal(err)
}
tok, err := conf.Exchange(ctx, code)

if err != nil {
    log.Fatal(err)
}
token =tok.AccessToken
writetoken(filePath,token)
}
//client := conf.Client(ctx, tok)
 config = dropbox.Config{
      Token: token,
      LogLevel: dropbox.LogOff, // if needed, set the desired logging level. Default is off
  }


}

func printFullAccount(w *tabwriter.Writer, fa *users.FullAccount) {
	fmt.Fprintf(w, "Logged in as %s <%s>\n\n", fa.Name.DisplayName, fa.Email)
	fmt.Fprintf(w, "Account Id:\t%s\n", fa.AccountId)

	fmt.Fprintf(w, "Account Id nuirmvjernvhunrhvn")
	fmt.Fprintf(w, "Account Type:\t%s\n", fa.AccountType.Tag)
	fmt.Fprintf(w, "Locale:\t%s\n", fa.Locale)
	fmt.Fprintf(w, "Referral Link:\t%s\n", fa.ReferralLink)
	fmt.Fprintf(w, "Profile Photo Url:\t%s\n", fa.ProfilePhotoUrl)
	fmt.Fprintf(w, "Paired Account:\t%t\n", fa.IsPaired)

}


func test() error{


	dbx := users.New(config)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 8, 1, ' ', 0)
	res, err := dbx.GetCurrentAccount()
		if err != nil {
			return err
		}
	fmt.Println(res.Name.DisplayName)
	
return w.Flush()

}


func uploadChunked(dbx files.Client, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64) (err error) {
	res, err := dbx.UploadSessionStart(files.NewUploadSessionStartArg(),
		&io.LimitedReader{R: r, N: chunkSize})
	if err != nil {
		return
	}

	written := chunkSize

	for (sizeTotal - written) > chunkSize {
		cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
		args := files.NewUploadSessionAppendArg(cursor)

		err = dbx.UploadSessionAppendV2(args, &io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return
		}
		written += chunkSize
	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	if _, err = dbx.UploadSessionFinish(args, r); err != nil {
		return
	}

	return
}



func put(src, dst string) (err error) {




	// Default `dst` to the base segment of the source path; use the second argument if provided.
	dst = "/" + path.Base(src)
	

	contents, err := os.Open(src)
	if err != nil {
		return
	}
	defer contents.Close()

	contentsInfo, err := contents.Stat()
	if err != nil {
		return
	}

	progressbar := &ioprogress.Reader{
		Reader: contents,
		DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
			return fmt.Sprintf("Uploading %s/%s",
				humanize.IBytes(uint64(progress)), humanize.IBytes(uint64(total)))
		}),
		Size: contentsInfo.Size(),
	}

	commitInfo := files.NewCommitInfo(dst)
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	commitInfo.ClientModified = time.Now().UTC().Round(time.Second)

	dbx := files.New(config)
	if contentsInfo.Size() > chunkSize {
		return uploadChunked(dbx, progressbar, commitInfo, contentsInfo.Size())
	}

	if _, err = dbx.Upload(commitInfo, progressbar); err != nil {
		return
	}

	return
}

func validatePath(p string) (path string, err error) {
	path = p

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	path = strings.TrimSuffix(path, "/")

	return
}

func mkdir(p string) (err error) {
	if p=="" {
		return errors.New("`mkdir` requires a `directory` argument")
	}

	dst, err := validatePath(p)
	if err != nil {
		return
	}

	arg := files.NewCreateFolderArg(dst)

	dbx := files.New(config)
	if _, err = dbx.CreateFolderV2(arg); err != nil {
		return
	}

	return
}

func main () {
confffiggg()
mkdir("pk")
fmt.Println(test())
var src string
fmt.Scan(&src);

put(src,"")
}
