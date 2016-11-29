# Chasm
Secure multi-party cloud backup solution based on Shamir's Secret Sharing scheme.

Presentation: https://alexgr.in/chasm

# Demo
[![asciicast](https://asciinema.org/a/2loda9ax8s22bvnl6nl5e728s.png)](https://asciinema.org/a/2loda9ax8s22bvnl6nl5e728s)

# Development
Make sure you have `godep` installed (`go get github.com/tools/godep`)

- Run `godep get` to install dependencies locally
- To build and execute: run `go build && ./chasm`
- Optionally, just run `godep go <CMD>` for any go command (ie: `godep go build && ./chasm`).

# Usage
```
NAME:
   chasm - A secret-sharing based secure cloud backup solution.
   
COMMANDS:
    start	Start running chasm.
    status	Prints out the current Chasm setup.
    add		Add a new cloud store to chasm.
    restore	Restores chasm after repeating setup.
    remove	Removes a cloud store.
    clean	Deletes all shares in cloud stores
    sync	Clean cloud stores, sync all items in Chasm folder by secret-sharing.
```
