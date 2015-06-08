/**
 *	(The MIT License)
 *
 *  Copyright (c) 2015 Waldemar Quevedo. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person
 *  obtaining a copy of this software and associated documentation
 *  files (the "Software"), to deal in the Software without
 *  restriction, including without limitation the rights to use, copy,
 *  modify, merge, publish, distribute, sublicense, and/or sell copies
 *  of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
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
		} else if os.Getenv("MESOS_SLAVE_PID") != "" {
			c.SetupEngine("mesos", data)
		} else {
			c.SetupEngine(*mode, data)
		}

		c.StartEngine()
		os.Exit(0)
	case "help", "h":
		cli.ShowUsage()
		os.Exit(0)
        }
}
