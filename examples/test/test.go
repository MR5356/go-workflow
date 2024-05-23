package main

import (
	"github.com/MR5356/go-workflow"
	"time"
)

type TestPlugin struct {
	workflow.UnimplementedITask
}

func (p *TestPlugin) SetParams(params *workflow.TaskParams) error {
	return nil
}

func (p *TestPlugin) Start() error {
	time.Sleep(time.Hour)
	return nil
}

func main() {
	workflow.Serve(&TestPlugin{})
}
