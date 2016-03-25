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
   0.0.1

COMMANDS:
   start	Start running chasm. start --root=<chasm_root>.
   add, a	Add a new cloud store to chasm. --root=<chasm_root> add <service>
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --root 			Chasm root directory. Example: --root=/home/alex
   --help, -h			show help
   --generate-bash-completion
   --version, -v		print the version
```
