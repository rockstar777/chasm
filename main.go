package main

// ChasmDir is the default directory for chasm
const ChasmDir = "~/Desktop/Chasm"

func main() {

	if CreateOrLoadChasmDir(ChasmDir) {
		RunSetup()
	}

}
