package main

import (
  "org-mode"
  "org-mode/export"
  "fmt"
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
