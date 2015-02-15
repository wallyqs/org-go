package main

import (
	"fmt"
	"io/ioutil"
	"org-mode"
	"org-mode/export"
	"os"
)

func main() {
	var org interface{}
	byteContents, _ := ioutil.ReadAll(os.Stdin)
	root := orgmode.Preprocess(string(byteContents))
	tokens := orgmode.Tokenize(string(byteContents), nil)

	org = orgmode.Parse(tokens, root)
	fmt.Println(orgexport.ToSexp(org))
}
