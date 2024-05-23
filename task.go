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
	t.client.Kill()
}

type ITask interface {
	DryRun() error
	GetWorkflow() *Workflow
	SetParams(params *TaskParams) error

	Start() error
	Stop() error
	Pause() error
	Resume() error
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

func (u *UnimplementedITask) DryRun() error {
	return nil
}

func (u *UnimplementedITask) GetWorkflow() *Workflow {
	return nil
}

func (u *UnimplementedITask) SetParams(params *TaskParams) error {
	return nil
}

func (u *UnimplementedITask) Start() error {
	return ErrNotImplemented
}

func (u *UnimplementedITask) Stop() error {
	return ErrNotImplemented
}

func (u *UnimplementedITask) Pause() error {
	return ErrNotImplemented
}

func (u *UnimplementedITask) Resume() error {
	return ErrNotImplemented
}
