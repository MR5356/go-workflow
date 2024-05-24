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

func UseLogger(name string) hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
}

var LoggerServer = hclog.New(&hclog.LoggerOptions{
	//Name:   "test",
	Output:     os.Stdout,
	JSONFormat: true,
	Level:      hclog.Info,
})
