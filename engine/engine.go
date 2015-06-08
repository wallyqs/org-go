package engine

import (
        "os/exec"
	"io"
	"bufio"
	"sync"
	"github.com/wallyqs/org-go/parser"
	"github.com/apcera/logray"
)

type CodeBlock struct {
	Name string
	Lang string
	Src  *orgmode.OrgSrcBlock
	Cmd *exec.Cmd
}

type Engine struct {
	Mode string
	CodeBlocks []*CodeBlock
	wg *sync.WaitGroup
	log *logray.Logger
}

func NewEngine(mode string, data []byte) *Engine {

        blocks := ProcessOrg(data)

	return &Engine {
		Mode: mode,
		CodeBlocks: blocks,
	        wg: &sync.WaitGroup{},
	}
}

// Takes a []byte with Org contents and filters
// for SrcBlock elements in the content
// FIXME: Should return error in case parsing fails
func ProcessOrg(data []byte) []*CodeBlock {
	root := orgmode.Preprocess(string(data))
	tokens := orgmode.Tokenize(string(data), root)

	blocks := make([]*CodeBlock, 0)

	for _, t := range tokens {
		switch o := t.(type) {
		case *orgmode.OrgSrcBlock:
		        // Create a code block
			
			var cmd *exec.Cmd
			switch o.Lang {
			case "ruby":
				cmd = exec.Command(o.Lang, "-e", o.RawContent)
			case "python":
				cmd = exec.Command(o.Lang, "-c", o.RawContent)
			case "sh":
				cmd = exec.Command(o.Lang, "-c", o.RawContent)
			case "js":
				// normalize to node and use the one in PATH
				cmd = exec.Command("node", "-e", o.RawContent)
			}

			block := &CodeBlock {
				Name: o.Name,
				Lang: o.Lang,
				Src:  o,
				Cmd:  cmd,
			}

			blocks = append(blocks, block)
		}
	}

	return blocks
}

func (e *Engine) Run () {

	switch e.Mode {
	case "local":
		// TODO: Only activate stdout when ':results output' in the future
		logray.AddDefaultOutput("stdout://", logray.ALL)
		e.log = logray.New()
		e.log.Infof("Starting Executor engine in '%s' mode...\n", e.Mode)		
		
		e.RunLocally()

	case "mesos":
		// TODO: Run locally as a Mesos executor
		e.log.Error("Not implemented yet...")
		break

	default:
		e.log.Error("Unrecognized mode:", e.Mode)
		break
	}
}

func (e *Engine) ExecuteCodeBlock(block *CodeBlock) (bool, error) {
	defer e.wg.Done()
	
        // goroutines, wait group, etc...
	r, w := io.Pipe()
	bufreader := bufio.NewReader(r)

	// Pipe both the stdout and stderr to a Reader
	block.Cmd.Stdout = w
	block.Cmd.Stderr = w

	e.log.Infof("Starting block: %v\n", block.Name)
	go func(bufr *bufio.Reader){
		logger := e.log.Clone()
		logger.SetField("block", block.Name)
		// FIXME: Need to Run first to get the process pid
		// logger.SetField("pid", block.Cmd.Process.Pid)
		for {
			// TODO: In case ReadString fails, then it should run a number of bytes from the buffer
			line, err := bufr.ReadString('\n')
			if err != nil {
				logger.Error("error while reading: ", err)
			}
			logger.Info(line)
			// fmt.Printf("[%d] %s -- %s", block.Cmd.Process.Pid, block.Name, line)
		}
	}(bufreader)

	// It is failing here???
	if err := block.Cmd.Run(); err != nil {
		e.log.Errorf("error during the execution of block: %v\n", err)
		return false, err
	}

	e.log.Infof("Block '%s' stopped running with: %v\n", block.Name, block.Cmd.ProcessState)
	return true, nil
}

func (e *Engine) RunLocally() {
	e.wg.Add(len(e.CodeBlocks))
	for _, block := range e.CodeBlocks {
		go e.ExecuteCodeBlock(block)
	}
	e.wg.Wait()
}
