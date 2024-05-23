package workflow

import (
	"github.com/hashicorp/go-hclog"
	"os"
)

var Logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Trace,
	Output:     os.Stderr,
	JSONFormat: true,
})

var LoggerServer = hclog.New(&hclog.LoggerOptions{
	//Name:   "test",
	Output: os.Stdout,
	Level:  hclog.Info,
})
