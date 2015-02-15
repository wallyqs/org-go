package main

import (
	"fmt"
	"org-mode"
	"org-mode/export"
)

func main() {

	example := `
  #+title: hello world

  #+options: todo:t ^:nil

  #+todo: TODO | DONE CANCELED

  Then any kind of text.

  `

	org := orgmode.Preprocess(example)
	fmt.Println(orgexport.ToSexp(org))
}
