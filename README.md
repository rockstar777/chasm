# chasm
Secure multi-party cloud backup solution based on Shamir's Secret Sharing scheme.

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
   add, a	Add a new cloud store to chasm.
   restore	Restores chasm after repeating setup.
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --root "/Users/Alex/Chasm"	Destination of the Chasm secure folder.
   --help, -h			show help
   --generate-bash-completion
   --version, -v		print the version
```
