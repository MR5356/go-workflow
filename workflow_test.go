package workflow

import (
	"reflect"
	"testing"
)

func TestWorkflow_AddTask(t *testing.T) {
	type fields struct {
		nodes []*Node
		edges []*Edge
	}
	type args struct {
		node *Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add node",
			fields: fields{
				nodes: []*Node{},
				edges: []*Edge{},
			},
			args: args{
				node: &Node{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				Nodes: tt.fields.nodes,
				Edges: tt.fields.edges,
			}
			w.AddNode(tt.args.node)
		})
	}
}

func TestWorkflow_AddEdge(t *testing.T) {
	type fields struct {
		nodes []*Node
		edges []*Edge
	}
	type args struct {
		source *Node
		target *Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add edge",
			fields: fields{
				nodes: []*Node{},
				edges: []*Edge{},
			},
			args: args{
				source: &Node{},
				target: &Node{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				Nodes: tt.fields.nodes,
				Edges: tt.fields.edges,
			}
			w.AddEdge(tt.args.source, tt.args.target)
		})
	}
}

type testTask struct {
	ID string
	UnimplementedITask
}

func TestWorkflow_HasCycle(t *testing.T) {
	t1 := &Node{Id: "1"}
	t2 := &Node{Id: "2"}
	t3 := &Node{Id: "3"}
	type fields struct {
		nodes []*Node
		edges []*Edge
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Node
		want1  bool
	}{
		{
			name: "has no cycle",
			fields: fields{
				nodes: []*Node{t1, t2, t3},
				edges: []*Edge{
					{Source: t1, Target: t2},
					{Source: t2, Target: t3},
				},
			},
			want:  nil,
			want1: false,
		},
		{
			name: "has cycle 1",
			fields: fields{
				nodes: []*Node{t1, t2, t3},
				edges: []*Edge{
					{Source: t1, Target: t2},
					{Source: t2, Target: t3},
					{Source: t3, Target: t1},
				},
			},
			want:  []*Node{t1, t2, t3, t1},
			want1: true,
		},
		{
			name: "has cycle 2",
			fields: fields{
				nodes: []*Node{t1, t2, t3},
				edges: []*Edge{
					{Source: t1, Target: t2},
					{Source: t2, Target: t1},
					{Source: t3, Target: t1},
				},
			},
			want:  []*Node{t1, t2, t1},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				Nodes: tt.fields.nodes,
				Edges: tt.fields.edges,
			}
			got, got1 := w.HasCycle()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HasCycle() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("HasCycle() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
