package workflow

import (
	"github.com/hashicorp/go-plugin"
	"os/exec"
)

const (
	NodeStatusPending = "pending"
	NodeStatusReady   = "ready"
	NodeStatusRunning = "running"
	NodeStatusFailure = "failure"
	NodeStatusSuccess = "success"
	NodeStatusAborted = "aborted"
)

func (n *Node) Dispense() (*Task, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         GetPluginMap(),
		// TODO 根据 n.Id 来获取插件及启动插件的方法
		Cmd:              exec.Command("sh", "-c", n.Uses),
		Logger:           LoggerServer,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	task := &Task{
		ITask:  &UnimplementedITask{},
		client: client,
	}

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	t, err := rpcClient.Dispense("task")
	if err != nil {
		return nil, err
	}
	var ok bool
	task.ITask, ok = t.(ITask)
	if !ok {
		return nil, ErrNotITask
	}
	return task, nil
}
