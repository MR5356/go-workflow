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

func (c *TaskGRPCClient) DryRun() error {
	_, err := c.client.DryRun(context.Background(), &Empty{})
	return err
}

func (s *TaskGRPCServer) DryRun(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, s.Impl.DryRun()
}

func (c *TaskGRPCClient) GetWorkflow() *Workflow {
	resp, err := c.client.GetWorkflow(context.Background(), &Empty{})
	if err != nil {
		return nil
	}
	return resp
}

func (s *TaskGRPCServer) GetWorkflow(ctx context.Context, req *Empty) (*Workflow, error) {
	return s.Impl.GetWorkflow(), nil
}

func (c *TaskGRPCClient) SetParams(params *TaskParams) error {
	_, err := c.client.SetParams(context.Background(), params)
	return err
}

func (s *TaskGRPCServer) SetParams(ctx context.Context, params *TaskParams) (*Empty, error) {
	return &Empty{}, s.Impl.SetParams(params)
}

func (c *TaskGRPCClient) Start() error {
	_, err := c.client.Start(context.Background(), &Empty{})
	return err
}

func (s *TaskGRPCServer) Start(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, s.Impl.Start()
}

func (c *TaskGRPCClient) Stop() error {
	_, err := c.client.Stop(context.Background(), &Empty{})
	return err
}

func (s *TaskGRPCServer) Stop(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, s.Impl.Stop()
}

func (c *TaskGRPCClient) Pause() error {
	_, err := c.client.Pause(context.Background(), &Empty{})
	return err
}

func (s *TaskGRPCServer) Pause(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, s.Impl.Pause()
}

func (c *TaskGRPCClient) Resume() error {
	_, err := c.client.Resume(context.Background(), &Empty{})
	return err
}

func (s *TaskGRPCServer) Resume(ctx context.Context, req *Empty) (*Empty, error) {
	return &Empty{}, s.Impl.Resume()
}
