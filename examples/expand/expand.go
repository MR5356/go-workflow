package main

import (
	"github.com/MR5356/go-workflow"
	"github.com/google/uuid"
)

type ExpandTask struct {
	workflow.UnimplementedITask
}

func (t *ExpandTask) GetWorkflow() *workflow.Workflow {
	n1 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "打工人",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n2 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "升职",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n4 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "加薪",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n3 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "躺平",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	wf := &workflow.Workflow{}
	wf.AddNode(n1, n2, n3, n4)
	wf.AddEdge(n1, n2)
	wf.AddEdge(n1, n3)
	wf.AddEdge(n2, n4)
	return wf
}

func main() {
	workflow.Serve(&ExpandTask{})
}
