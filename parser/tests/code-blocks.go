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
	orgTestFilePath := ospath.Join(currentDir, "org/features/complex/src-1.org")

	contents, err := io.ReadFile(orgTestFilePath)
	if err != nil {
		fmt.Printf("Problem reading the file: %v \n", err)
	}

	var org interface{}
	root := orgmode.Preprocess(string(contents))
	tokens := orgmode.Tokenize(string(contents), root)

	org = orgmode.Parse(tokens, root)
	fmt.Println(orgexport.ToSexp(org))
}
