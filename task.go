package workflow

import (
	"errors"
	"github.com/hashicorp/go-plugin"
)

var (
	ErrNotImplemented = errors.New("method not implemented")
	ErrNotITask       = errors.New("not an ITask")
)

type Task struct {
	ITask
	client *plugin.Client
}

func (t *Task) Close() {
	if t != nil && t.client != nil {
		t.client.Kill()
	}
}

type ITask interface {
	GetWorkflow() *WorkflowDAG
	SetParams(params *TaskParams) error

	Run() error
}

func (ps *TaskParams) Get(key string) string {
	for _, p := range ps.Params {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

type UnimplementedITask struct{}

func (u *UnimplementedITask) GetWorkflow() *WorkflowDAG {
	return nil
}

func (u *UnimplementedITask) SetParams(params *TaskParams) error {
	return nil
}

func (u *UnimplementedITask) Run() error {
	return ErrNotImplemented
}
