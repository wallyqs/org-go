package main

import (
  . "org-mode"
    "fmt"
    "strconv"
    "os/exec"
    "bytes"
)

func main() {

  raw := `

Running code blocks sequentially.

First a Ruby block:

#+name: ruby-hello-world
#+BEGIN_SRC ruby
(0..20).each do |n|
  puts "#{n}: Hello World from Ruby!"
end
#+END_SRC

and then a Python block:

#+name: python-hello-world
#+BEGIN_SRC python
print "hello world from python!!!!"
#+END_SRC

`

  fmt.Println("----- Org Content ---")
  fmt.Println(raw)
  fmt.Println("---------------------")

  root   := Preprocess(raw)
  tokens := Tokenize(raw, root)

  // Simple filter
  blocks := make([]*OrgSrcBlock, 0)
  for _, t := range tokens {
    switch o := t.(type) {
      case *OrgSrcBlock:
        blocks = append(blocks, o)
    }
  }

  for _, codeblock := range blocks {
    fmt.Println("■■■ Code Block Info ■■■")
    fmt.Println("Name:    ", codeblock.Name)
    fmt.Println("Lang:    ", codeblock.Lang)
    fmt.Println("Value:   ", strconv.Quote(codeblock.RawContent))
    fmt.Println("Headers: ", codeblock.Headers, "\n")

    // Now let's try to execute it

    var cmd *exec.Cmd;
    switch codeblock.Lang {
      case "ruby":
        cmd = exec.Command(codeblock.Lang, "-e", codeblock.RawContent)
      case "python":
        cmd = exec.Command(codeblock.Lang, "-c", codeblock.RawContent)
    }

    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()

    fmt.Println(",#+RESULTS:", codeblock.Name)
    if err != nil {
      fmt.Println("Execution Failed: ", err)
    }

    if len(stdout.String()) > 0 {
      fmt.Println(",#+begin_example\n")
      fmt.Println(stdout.String())
      fmt.Println(",#+end_example")
    } else {
      fmt.Println("No output")
    }

    if len(stderr.String()) > 0 {
      fmt.Println("\n,#+begin_error:\n")
      fmt.Println(stderr.String())
      fmt.Println(",#+end_error")
    }

  }
}
