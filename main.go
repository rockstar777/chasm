package main

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"os"
)

/// chasm commands ///

func loadChasm(c *cli.Context) {
	if chasmRoot == "" {
		color.Red("Error: missing chasm root path.")
		os.Exit(2)
		return
	}

	CreateOrLoadChasmDir(chasmRoot)
}

func startChasm(c *cli.Context) {
	loadChasm(c)

	if preferences.NeedSetup() {
		color.Red("Error: not enough services. Add a service with add .")
		return
	}

	// start the watcher
	StartWatching(preferences.root)
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

func addDropbox(c *cli.Context) {
	loadChasm(c)
	color.Red("Error: not implemented.")
}

func addDrive(c *cli.Context) {
	loadChasm(c)
	color.Red("Error: not implemented.")
}


/// Cli toolchain ///
var chasmRoot string

func main() {
	app := cli.NewApp()

	app.Name = "chasm"
	app.Usage = "A secret-sharing based secure cloud backup solution."
	app.EnableBashCompletion = true
	app.Version = "0.0.1"

	app.Flags = []cli.Flag {
	  cli.StringFlag{
	    Name: "root",
		Value: "",
	    Usage: "Chasm root directory. Example: --root=/home/alex",
		Destination: &chasmRoot,
	  },
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: nil,
			Usage:   "Start running chasm. start --root=<chasm_root>.",
			Action: startChasm,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new cloud store to chasm. --root=<chasm_root> add <service>",
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

	}

	app.Run(os.Args)
}
