package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

/// chasm commands ///

func loadChasm(c *cli.Context) {
	CreateOrLoadChasmDir(chasmRoot)
}

func startChasm(c *cli.Context) {
	loadChasm(c)

	if preferences.NeedSetup() {
		color.Red("Warning: not enough services.")
		return
	}

	// start the watcher
	color.Green("Starting chasm. Listening on %s", preferences.root)
	StartWatching(preferences.root)
}

func statusChasm(c *cli.Context) {
	loadChasm(c)

	color.Green("Cloud stores:")
	for i, cs := range preferences.AllCloudStores() {
		fmt.Println(color.GreenString("%v)", i+1), cs.Description())
	}
	if preferences.NeedSetup() {
		color.Red("Warning: not enough services.")
	}
}

func restoreChasm(c *cli.Context) {
	loadChasm(c)

	if preferences.NeedSetup() {
		color.Red("Warning: not enough services. Cannot Restore.")
		return
	}

	color.Green("Preparing to restore chasm to %s", preferences.root)
	Restore()
}

func cleanChasm(c *cli.Context) {
	loadChasm(c)

	for _, cs := range preferences.AllCloudStores() {
		cs.Clean()
	}
}

func addFolder(c *cli.Context) {
	loadChasm(c)
	var folderStore FolderStore

	if len(c.Args()) < 1 {
		color.Red("Error: missing folder path")
		return
	}

	folderStore.Path = c.Args()[0]
	folderStore.Setup()
	preferences.FolderStores = append(preferences.FolderStores, folderStore)
	preferences.Save()

	color.Green("Success! Added folder store: %s", folderStore.Path)
}

func addDrive(c *cli.Context) {
	loadChasm(c)
	var gdrive GDriveStore

	if (&gdrive).Setup() == false {
		color.Red("(Cloud Store) Google Drive: setup incomplete.")
		return
	}

	// only 1 gdrive store
	preferences.GDriveStores = []GDriveStore{gdrive}
	preferences.Save()

	color.Green("Success! Added Google Drive Store.")
}

func addDropbox(c *cli.Context) {
	loadChasm(c)

	var dropbox DropboxStore

	if (&dropbox).Setup() == false {
		color.Red("(Cloud Store) Dropbox: setup incomplete.")
		return
	}

	// only 1 dropbox store
	preferences.DropboxStores = []DropboxStore{dropbox}
	preferences.Save()

	color.Green("Success! Added Dropbox Store.")
}

/// Cli toolchain ///
var chasmRoot string

func main() {
	app := cli.NewApp()

	app.Name = color.GreenString("chasm")
	app.Usage = color.GreenString("A secret-sharing based secure cloud backup solution.")
	app.EnableBashCompletion = true
	app.Version = "0.1"

	usr, _ := user.Current()
	defaultRoot := path.Join(usr.HomeDir, "Chasm")
	chasmRoot = defaultRoot

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "root",
			Value:       defaultRoot,
			Usage:       "Destination of the Chasm secure folder.",
			Destination: &chasmRoot,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: nil,
			Usage:   "Start running chasm.",
			Action:  startChasm,
		},
		{
			Name:    "status",
			Aliases: nil,
			Usage:   "Prints out the current Chasm setup.",
			Action:  statusChasm,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new cloud store to chasm.",
			Subcommands: []cli.Command{
				{
					Name:   "folder",
					Usage:  "add folder",
					Action: addFolder,
				},
				{
					Name:   "dropbox",
					Usage:  "add dropbox",
					Action: addDropbox,
				},
				{
					Name:   "drive",
					Usage:  "add google drive",
					Action: addDrive,
				},
			},
		},
		{
			Name:    "restore",
			Aliases: nil,
			Usage:   "Restores chasm after repeating setup.",
			Action:  restoreChasm,
		},
		{
			Name:    "clean",
			Aliases: nil,
			Usage:   "Deletes all shares in cloud stores",
			Action:  cleanChasm,
		},
	}

	app.Run(os.Args)
}
