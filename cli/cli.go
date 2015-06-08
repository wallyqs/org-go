package cli

import (
	"fmt"
	"strings"
	"os"
	"github.com/wallyqs/org-go/engine"
	"github.com/apcera/termtables"
)

type CLI struct {
	engine *engine.Engine
}

// Takes options and returns a CLI with an attached Engine
func NewCLI(args []string) *CLI {
	c := &CLI{}
	return c
}

func ShowExecution (blocks []*engine.CodeBlock) {
	table := termtables.CreateTable()
	table.UTF8Box()
	table.AddHeaders("Name", "Lang", "Headers")

	for _, block := range blocks {
		var headers string
		for k, v := range block.Src.Headers {
			headers += strings.Join([]string{k, v}, " ") + " "
		}

		table.AddRow(block.Name, block.Lang, headers)
	}

	fmt.Println(table.Render())
}

func ShowUsage () {
	usage := `
  org-go [-m mode] COMMAND [OPTIONS]

  Available commands:

  run:   Execute the code blocks within the Org mode file
  show:  Summary of code blocks that will be executed
  help:  Show this message

`
	fmt.Fprint(os.Stderr, usage)
}

func (c *CLI) SetupEngine(mode string, data []byte) {
        c.engine = engine.NewEngine(mode, data)
}

func (c *CLI) StartEngine() {
	c.engine.Run()
}
