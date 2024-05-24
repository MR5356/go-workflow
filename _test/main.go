package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MR5356/go-workflow"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

var w *workflow.Workflow
var taskQueue = make(chan *workflow.Node, 10)
var ctx, cancel = context.WithCancel(context.Background())

func Runner() {
	for node := range taskQueue {
		logrus.Infof("run node: %+v", node)
		err := w.RunNode(node)
		if err != nil {
			logrus.Errorf("run node error: %+v", err)
		}
		logrus.Infof("run node done: %+v", node)
	}
}

func getWF() *workflow.Workflow {
	n1 := &workflow.Node{
		Id:    uuid.NewString(),
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
		Id:    uuid.NewString(),
		Label: "buildx环境准备",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n3 := &workflow.Node{
		Id:    uuid.NewString(),
		Uses:  "./_output/plugin/sleep/sleep",
		Label: "镜像仓库认证",
	}
	n41 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n42 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n43 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n44 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n51 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n52 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n53 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n54 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n6 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "部署PreProd",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n71 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "冒烟测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n72 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "API测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n73 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "性能测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n8 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "部署Prod",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	wf := workflow.NewWorkflow(ctx, taskQueue)
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
		Id:    uuid.NewString(),
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
		Id:    uuid.NewString(),
		Label: "buildx环境准备",
		Uses:  "./_output/plugin/expand/expand",
	}
	n3 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像仓库认证",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n41 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n42 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n43 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n44 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "编译4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n51 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描1",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n52 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描2",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n53 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描3",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n54 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "镜像扫描4",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n6 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "部署PreProd",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n71 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "冒烟测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n72 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "API测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n73 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "性能测试",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	n8 := &workflow.Node{
		Id:    uuid.NewString(),
		Label: "部署Prod",
		Uses:  "./_output/plugin/sleep/sleep",
	}
	wf := workflow.NewWorkflow(ctx, taskQueue)
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
	for i := 0; i < 3; i++ {
		go Runner()
	}
	w = getWF()

	http.HandleFunc("/api/v1", func(writer http.ResponseWriter, request *http.Request) {
		ww := w.GetWorkflow()
		s, _ := json.Marshal(ww)

		fmt.Fprintf(writer, string(s))
	})

	http.HandleFunc("/api/v1/start", func(writer http.ResponseWriter, request *http.Request) {
		err := w.Run()
		if err != nil {
			fmt.Fprintf(writer, string("failed"))
		}

		fmt.Fprintf(writer, string("success"))
	})

	http.HandleFunc("/api/v1/cancel", func(writer http.ResponseWriter, request *http.Request) {
		cancel()
		fmt.Fprintf(writer, string("success"))
	})

	http.HandleFunc("/api/v1/reset", func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel = context.WithCancel(context.Background())
		w = getExpandWF()
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Fatalf("listen error: %+v", err)
	}
}
