package workflow

import "github.com/hashicorp/go-plugin"

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TASK_PLUGIN",
	MagicCookieValue: "WORKFLOW",
}

func Serve(task ITask) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         GetPluginMap(task),
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

func GetPluginMap(task ...ITask) map[string]plugin.Plugin {
	if len(task) > 0 {
		return map[string]plugin.Plugin{
			"task": &TaskGRPCPlugin{Impl: task[0]},
		}
	} else {
		return map[string]plugin.Plugin{
			"task": &TaskGRPCPlugin{},
		}
	}
}
