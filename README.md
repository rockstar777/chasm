# chasm
Secure multi-party cloud backup solution based on Shamir's Secret Sharing scheme.

Demo:
[![asciicast](https://asciinema.org/a/2loda9ax8s22bvnl6nl5e728s.png)](https://asciinema.org/a/2loda9ax8s22bvnl6nl5e728s)

# development
Make sure you have `godep` installed (`go get github.com/tools/godep`)

- Run `godep get` to install dependencies locally
- To build and execute: run `go build && ./chasm`
- Optionally, just run `godep go <CMD>` for any go command (ie: `godep go build && ./chasm`).

# usage
```
NAME:
   chasm - A secret-sharing based secure cloud backup solution.

USAGE:
   chasm [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
    start	Start running chasm.
    status	Prints out the current Chasm setup.
    add		Add a new cloud store to chasm.
    restore	Restores chasm after repeating setup.
    remove	Removes a cloud store.
    clean	Deletes all shares in cloud stores
    sync	Clean cloud stores, sync all items in Chasm folder by secret-sharing.

GLOBAL OPTIONS:
   --root value, -r value	Destination of the Chasm secure folder. (default: "/Users/Alex/Chasm")
   --help, -h			show help
   --version, -v		print the version
```
