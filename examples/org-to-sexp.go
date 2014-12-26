package main

import (
	"org-mode"
	"org-mode/export"
	"io/ioutil"
	"fmt"
	"os"
)

func main() {
	var org interface{}
	byteContents, _ := ioutil.ReadAll(os.Stdin)
	root   := orgmode.Preprocess(string(byteContents))
	tokens := orgmode.Tokenize(string(byteContents), nil)

	org = orgmode.Parse(tokens, root)
	fmt.Println(orgexport.ToSexp(org))
}
