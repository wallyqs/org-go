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
	orgTestFilePath := ospath.Join(currentDir, "org/features/headlines/levels.org")

	contents, err := io.ReadFile(orgTestFilePath)
	if err != nil {
		fmt.Printf("Problem reading the file: %v \n", err)
	}

	// Still return an Org mode root which defines the current context
	// TODO: preprocessor reference does not work...
	root := orgmode.Preprocess(string(contents))

	// return an []interface{} of tokens
	tokens := orgmode.Tokenize(string(contents), root)

	var org interface{}
	org = orgmode.Parse(tokens, root)
	fmt.Println(orgexport.ToSexp(org))
}
