package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	. "org-mode"
	"os"
	"os/exec"
)

func main() {

	rawbytes, _ := ioutil.ReadAll(os.Stdin)
	raw := string(rawbytes)
	root := Preprocess(string(raw))
	tokens := Tokenize(string(raw), root)

	// Simple filter
	blocks := make([]*OrgSrcBlock, 0)
	for _, t := range tokens {
		switch o := t.(type) {
		case *OrgSrcBlock:
			blocks = append(blocks, o)
		}
	}

	for i, codeblock := range blocks {
		fmt.Println("\n* Code Block: ", i, "\n\n")

		fmt.Println("#+name: ", codeblock.Name)
		fmt.Println("#+begin_src ", codeblock.Lang, " ", codeblock.Headers)
		fmt.Println(codeblock.RawContent)
		fmt.Println("#+end_src", "\n")

		// Now let's try to execute it

		var cmd *exec.Cmd
		switch codeblock.Lang {
		case "ruby":
			cmd = exec.Command(codeblock.Lang, "-e", codeblock.RawContent)
		case "python":
			cmd = exec.Command(codeblock.Lang, "-c", codeblock.RawContent)
		case "sh":
			cmd = exec.Command(codeblock.Lang, "-c", codeblock.RawContent)
		case "js":
			cmd = exec.Command(codeblock.Lang, "-e", codeblock.RawContent)
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()

		fmt.Println("#+RESULTS:", codeblock.Name)
		if err != nil {
			fmt.Println(": Execution Failed: ", err)
		}

		if len(stdout.String()) > 0 {
			fmt.Println("#+begin_example\n")
			fmt.Println(stdout.String())
			fmt.Println("#+end_example")
		} else {
			fmt.Println("No output")
		}

		if len(stderr.String()) > 0 {
			fmt.Println("\n#+BEGIN_ERROR\n")
			fmt.Println(stderr.String())
			fmt.Println("#+END_ERROR")
		}

	}
}
