package main

import (
	"fmt"
	io "io/ioutil"
	"org-mode"
	"org-mode/export"
	"os"
	ospath "path/filepath"
)

func main() {
	currentDir, _ := os.Getwd()
	orgTestFilePath := ospath.Join(currentDir, "org/features/options/title.org")

	contents, err := io.ReadFile(orgTestFilePath)
	if err != nil {
		fmt.Printf("Problem reading the file: %v \n", err)
	}

	// TODO: process the string and return the first Org root
	org := orgmode.Preprocess(string(contents))

	fmt.Println(orgexport.ToSexp(org))
}
