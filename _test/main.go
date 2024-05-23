package main

import (
	"encoding/json"
	"fmt"
	"github.com/MR5356/go-workflow"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
)

var we *WorkflowExecutor

type WorkflowExecutor struct {
	*workflow.Workflow
	taskQueue chan *workflow.Node
}

func NewWorkflowExecutor(wf *workflow.Workflow) *WorkflowExecutor {
	return &WorkflowExecutor{
		Workflow:  wf,
		taskQueue: make(chan *workflow.Node, 10),
	}
}

func (we *WorkflowExecutor) DryRun() {
	for _, task := range we.Nodes {
		tk, err := task.Dispense()
		defer tk.Close()
		if err != nil {
			logrus.Errorf("dispense error: %+v", err)
			continue
		}
		err = tk.SetParams(&workflow.TaskParams{Params: task.Params})
		if err != nil {
			logrus.Errorf("set params error: %+v", err)
			continue
		}
		subWf := tk.GetWorkflow()
		if !subWf.IsNil() {
			we.ReplaceNodeToSubWorkflow(task, subWf)
		}
	}
	logrus.Infof("workflow: %+v", we.Nodes)
}

func (we *WorkflowExecutor) Runner() {
	for task := range we.taskQueue {
		task.Status = workflow.NodeStatusRunning
		tk, err := task.Dispense()
		defer tk.Close()
		if err != nil {
			logrus.Errorf("dispense error: %+v", err)
			task.Status = workflow.NodeStatusFailure
			continue
		}
		err = tk.SetParams(&workflow.TaskParams{Params: task.Params})
		if err != nil {
			logrus.Errorf("set params error: %+v", err)
			task.Status = workflow.NodeStatusFailure
			continue
		}
		err = tk.Start()
		if err != nil {
			logrus.Errorf("start error: %+v", err)
			task.Status = workflow.NodeStatusFailure
			continue
		}
		task.Status = workflow.NodeStatusSuccess
		ns := we.GetReadyNodes()
		for _, n := range ns {
			we.taskQueue <- n
		}
	}
}

func (we *WorkflowExecutor) Executor() {
	for i := 0; i < 2; i++ {
		go we.Runner()
	}
	we.DryRun()
	ns := we.GetReadyNodes()
	for _, n := range ns {
		n.Status = workflow.NodeStatusReady
		we.taskQueue <- n
	}
}

func getWF() *workflow.Workflow {
	n1 := &workflow.Node{
		Id:    "node1",
		Label: "克隆代码",
		Uses:  "./_output/plugin/sleep/sleep",
		Params: []*workflow.TaskParam{
			{
				Key:   "sleep",
				Value: "5",
			},
		},
	}
	n2 := &workflow.Node{
		Id:    "node2",
		Label: "buildx环境准备",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n3 := &workflow.Node{
		Id:    "node3",
		Uses:  "./_output/plugin/sleep/sleep",
		Label: "镜像仓库认证",
	}
	n41 := &workflow.Node{
		Id:    "node41",
		Label: "编译1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n42 := &workflow.Node{
		Id:    "node42",
		Label: "编译2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n43 := &workflow.Node{
		Id:    "node43",
		Label: "编译3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n44 := &workflow.Node{
		Id:    "node44",
		Label: "编译4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n51 := &workflow.Node{
		Id:    "node51",
		Label: "镜像扫描1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n52 := &workflow.Node{
		Id:    "node52",
		Label: "镜像扫描2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n53 := &workflow.Node{
		Id:    "node53",
		Label: "镜像扫描3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n54 := &workflow.Node{
		Id:    "node54",
		Label: "镜像扫描4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n6 := &workflow.Node{
		Id:    "node6",
		Label: "部署PreProd",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n71 := &workflow.Node{
		Id:    "node71",
		Label: "冒烟测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n72 := &workflow.Node{
		Id:    "node72",
		Label: "API测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n73 := &workflow.Node{
		Id:    "node73",
		Label: "性能测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n8 := &workflow.Node{
		Id:    "node8",
		Label: "部署Prod",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	wf := &workflow.Workflow{}
	wf.AddNode(n1, n2, n3, n41, n42, n43, n44, n51, n52, n53, n54, n6, n71, n72, n73, n8)
	wf.AddEdge(n1, n2)
	wf.AddEdge(n2, n3)
	wf.AddEdge(n3, n41)
	wf.AddEdge(n3, n42)
	wf.AddEdge(n3, n43)
	wf.AddEdge(n3, n44)
	wf.AddEdge(n41, n51)
	wf.AddEdge(n42, n52)
	wf.AddEdge(n43, n53)
	wf.AddEdge(n44, n54)
	wf.AddEdge(n51, n6)
	wf.AddEdge(n52, n6)
	wf.AddEdge(n53, n6)
	wf.AddEdge(n54, n6)
	wf.AddEdge(n6, n71)
	wf.AddEdge(n6, n72)
	wf.AddEdge(n6, n73)
	wf.AddEdge(n71, n8)
	wf.AddEdge(n72, n8)
	//wf.AddEdge(n73, n8)
	return wf
}

func getExpandWF() *workflow.Workflow {
	n1 := &workflow.Node{
		Id:    "node1",
		Label: "克隆代码",
		Uses:  "./_output/plugin/sleep/sleep",
		Params: []*workflow.TaskParam{
			{
				Key:   "sleep",
				Value: "5",
			},
		},
	}
	n2 := &workflow.Node{
		Id:    "node2",
		Label: "buildx环境准备",
		Uses:  "./_output/plugin/expand/expand",
	}
	n3 := &workflow.Node{
		Id:    "node3",
		Label: "镜像仓库认证",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n41 := &workflow.Node{
		Id:    "node41",
		Label: "编译1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n42 := &workflow.Node{
		Id:    "node42",
		Label: "编译2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n43 := &workflow.Node{
		Id:    "node43",
		Label: "编译3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n44 := &workflow.Node{
		Id:    "node44",
		Label: "编译4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n51 := &workflow.Node{
		Id:    "node51",
		Label: "镜像扫描1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n52 := &workflow.Node{
		Id:    "node52",
		Label: "镜像扫描2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n53 := &workflow.Node{
		Id:    "node53",
		Label: "镜像扫描3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n54 := &workflow.Node{
		Id:    "node54",
		Label: "镜像扫描4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n6 := &workflow.Node{
		Id:    "node6",
		Label: "部署PreProd",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n71 := &workflow.Node{
		Id:    "node71",
		Label: "冒烟测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n72 := &workflow.Node{
		Id:    "node72",
		Label: "API测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n73 := &workflow.Node{
		Id:    "node73",
		Label: "性能测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n8 := &workflow.Node{
		Id:    "node8",
		Label: "部署Prod",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	wf := &workflow.Workflow{}
	wf.AddNode(n1, n2, n3, n41, n42, n43, n44, n51, n52, n53, n54, n6, n71, n72, n73, n8)
	wf.AddEdge(n1, n2)
	wf.AddEdge(n2, n3)
	wf.AddEdge(n3, n41)
	wf.AddEdge(n3, n42)
	wf.AddEdge(n3, n43)
	wf.AddEdge(n3, n44)
	wf.AddEdge(n41, n51)
	wf.AddEdge(n42, n52)
	wf.AddEdge(n43, n53)
	wf.AddEdge(n44, n54)
	wf.AddEdge(n51, n6)
	wf.AddEdge(n52, n6)
	wf.AddEdge(n53, n6)
	wf.AddEdge(n54, n6)
	wf.AddEdge(n6, n71)
	wf.AddEdge(n6, n72)
	wf.AddEdge(n6, n73)
	wf.AddEdge(n71, n8)
	wf.AddEdge(n72, n8)
	//wf.AddEdge(n73, n8)
	return wf
}

func main() {
	we = NewWorkflowExecutor(getWF())

	http.HandleFunc("/api/v1", func(writer http.ResponseWriter, request *http.Request) {
		ww := we.GetWorkflow()
		s, _ := json.Marshal(ww)

		fmt.Fprintf(writer, string(s))
	})

	http.HandleFunc("/api/v1/start", func(writer http.ResponseWriter, request *http.Request) {
		we.Executor()
		fmt.Fprintf(writer, string("success"))
	})

	http.HandleFunc("/api/v1/reset", func(writer http.ResponseWriter, request *http.Request) {
		we = NewWorkflowExecutor(getExpandWF())
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Fatalf("listen error: %+v", err)
	}
}

func mainPlugin() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  workflow.Handshake,
		Plugins:          workflow.GetPluginMap(),
		Cmd:              exec.Command("sh", "-c", "./_output/plugin/checkout"),
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		logrus.Fatalf("client.Client() failed: %+v", err)
	}

	task, err := rpcClient.Dispense("task")
	if err != nil {
		logrus.Fatalf("rpcClient.Dispense() failed: %+v", err)
	}

	err = task.(workflow.ITask).Start()
	if err != nil {
		logrus.Fatalf("task.Start() failed: %+v", err)
	}
}
