package workflow

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type TaskGRPCPlugin struct {
	plugin.Plugin
	Impl ITask
}

func (p *TaskGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterTaskServer(s, &TaskGRPCServer{Impl: p.Impl})
	return nil
}

func (p *TaskGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &TaskGRPCClient{client: NewTaskClient(c)}, nil
}

type TaskGRPCClient struct {
	client TaskClient
}

type TaskGRPCServer struct {
	Impl ITask
	UnimplementedTaskServer
}

func (c *TaskGRPCClient) GetWorkflow() *WorkflowDAG {
	resp, err := c.client.GetWorkflow(context.Background(), &Empty{})
	if err != nil {
		return nil
	}
	return resp
}

func (s *TaskGRPCServer) GetWorkflow(ctx context.Context, req *Empty) (*WorkflowDAG, error) {
	return s.Impl.GetWorkflow(), nil
}

func (c *TaskGRPCClient) SetParams(params *TaskParams) error {
	_, err := c.client.SetParams(context.Background(), params)
	return err
}

func (s *TaskGRPCServer) SetParams(ctx context.Context, params *TaskParams) (*Empty, error) {
	return &Empty{}, s.Impl.SetParams(params)
}

func (c *TaskGRPCClient) Run() error {
	_, err := c.client.Run(context.Background(), &Empty{})
	return err
}

func (s *TaskGRPCServer) Run(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, s.Impl.Run()
}
