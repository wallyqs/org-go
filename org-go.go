package main

import (
        "fmt"
        "flag"
        "os"
        "io/ioutil"
        "github.com/wallyqs/org-go/cli"
	"github.com/wallyqs/org-go/engine"
)

var mode = flag.String("m", "local", "Running mode for the executor")

func init() {
        flag.Parse()
}

func main() {
        args := flag.Args()

        c := cli.NewCLI(args)

        if len(args) < 2 {
                cli.ShowUsage()
                os.Exit(1)
        }

        // check for command used
        command := args[0]

        data, err := ioutil.ReadFile(args[1])
        if err != nil {
                cli.ShowUsage()
                fmt.Fprintf(os.Stderr, "error reading Org file: %v\n", err)
                os.Exit(1)
        }

        switch command {
        case "show", "s":
	        blocks := engine.ProcessOrg(data)
                cli.ShowExecution(blocks)
                os.Exit(0)
        case "run", "r":
	        if *mode == "" {
		        c.SetupEngine("local", data)
			c.StartEngine()
		} else {
			c.SetupEngine(*mode, data)
			c.StartEngine()
		}

		os.Exit(0)
	case "help", "h":
		cli.ShowUsage()
		os.Exit(0)
        }
}
