package workflow

import (
	"context"
)

type Workflow struct {
	*WorkflowDAG
	nodeMap   map[string]*Node
	taskQueue chan *Node
	ctx       context.Context
	stop      bool
}

func NewWorkflow(ctx context.Context, tq chan *Node) *Workflow {
	return &Workflow{
		WorkflowDAG: &WorkflowDAG{},
		nodeMap:     make(map[string]*Node),
		taskQueue:   tq,
		ctx:         ctx,
		stop:        false,
	}
}

func (w *WorkflowDAG) AddNode(node ...*Node) {
	for _, n := range node {
		n.Status = NodeStatusPending
	}
	w.Nodes = append(w.Nodes, node...)
}

func (w *WorkflowDAG) AddEdge(source, target *Node) {
	w.Edges = append(w.Edges, &Edge{Source: source.Id, Target: target.Id})
}

func (w *Workflow) AddNode(node ...*Node) {
	for _, n := range node {
		n.Status = NodeStatusPending
		w.nodeMap[n.Id] = n
	}
	w.Nodes = append(w.Nodes, node...)
}

func (w *Workflow) AddEdge(source, target *Node) {
	w.Edges = append(w.Edges, &Edge{Source: source.Id, Target: target.Id})
}

func (w *Workflow) getSubWorkflow(node *Node) (*WorkflowDAG, error) {
	task, err := node.Dispense()
	defer task.Close()
	if err != nil {
		return nil, err
	}
	err = task.SetParams(&TaskParams{Params: node.Params})
	if err != nil {
		return nil, err
	}

	return task.GetWorkflow(), nil
}

func (w *Workflow) ExpandDynamicNode() error {
	for _, node := range w.Nodes {
		subWf, err := w.getSubWorkflow(node)
		if err != nil {
			return err
		}
		if !subWf.IsNil() {
			w.ReplaceNodeToSubWorkflow(node, subWf)
		}
	}

	return nil
}

func (w *Workflow) RunNode(node *Node) error {
	node.Status = NodeStatusRunning
	task, err := node.Dispense()
	defer task.Close()
	if err != nil {
		node.Status = NodeStatusFailure
		return err
	}

	errChan := make(chan error)
	go func() {
		err = task.SetParams(&TaskParams{Params: node.Params})
		if err != nil {
			node.Status = NodeStatusFailure
			errChan <- err
		}

		err = task.Run()
		if err != nil {
			node.Status = NodeStatusFailure
			errChan <- err
		}

		node.Status = NodeStatusSuccess
		w.Next()
		errChan <- nil
	}()
	select {
	case err = <-errChan:
		return err
	case <-w.ctx.Done():
		task.Close()
		w.stop = true
		return w.ctx.Err()
	}
}

func (w *Workflow) Next() {
	if w.stop {
		return
	}
	ns := w.GetReadyNodes()
	for _, n := range ns {
		//logrus.Infof("run node ready: %s", n.Id)
		n.Status = NodeStatusReady
		w.taskQueue <- n
	}
}

func (w *Workflow) Run() error {
	err := w.ExpandDynamicNode()
	if err != nil {
		return err
	}
	w.Next()
	return nil
}

func (w *WorkflowDAG) IsNil() bool {
	return w == nil || len(w.Nodes) == 0
}

func (w *Workflow) ReplaceNodeToSubWorkflow(source *Node, target *WorkflowDAG) {
	// 替换节点
	for _, n := range target.Nodes {
		w.nodeMap[n.Id] = n
	}
	w.Nodes = append(w.Nodes, target.Nodes...)
	w.Edges = append(w.Edges, target.Edges...)
	for i, node := range w.Nodes {
		if node == source {
			w.Nodes = append(w.Nodes[:i], w.Nodes[i+1:]...)
		}
	}

	// 替换边
	temps := make(map[string]int)
	tempt := make(map[string]int)
	for _, t := range target.Nodes {
		temps[t.Id] = 0
		tempt[t.Id] = 0
	}

	for _, e := range target.Edges {
		temps[e.Source]++
		tempt[e.Target]++
	}

	for index := len(w.Edges) - 1; index >= 0; index-- {
		edge := w.Edges[index]
		if edge.Source == source.Id {
			var ts []*Node
			for t, degree := range temps {
				if degree == 0 {
					ts = append(ts, w.nodeMap[t])
				}
			}
			for _, t := range ts {
				w.Edges = append(w.Edges, &Edge{Source: t.Id, Target: edge.Target})
			}

			// 删除边
			w.Edges = append(w.Edges[:index], w.Edges[index+1:]...)
		}
		if edge.Target == source.Id {
			var ts []*Node
			for t, degree := range tempt {
				if degree == 0 {
					ts = append(ts, w.nodeMap[t])
				}
			}

			for _, t := range ts {
				w.Edges = append(w.Edges, &Edge{Source: edge.Source, Target: t.Id})
			}

			// 删除边
			w.Edges = append(w.Edges[:index], w.Edges[index+1:]...)
		}
	}
}

func (w *Workflow) GetWorkflow() *WorkflowDAG {
	for _, edge := range w.Edges {
		edge.Status = w.nodeMap[edge.Target].Status
	}
	return w.WorkflowDAG
}

func (w *WorkflowDAG) isFinalState(node *Node) bool {
	return node.Status == NodeStatusSuccess || node.Status == NodeStatusFailure || node.Status == NodeStatusAborted
}

func (w *Workflow) GetReadyNodes() []*Node {
	inDegree := make(map[string]int)
	for _, node := range w.Nodes {
		if node.Status == NodeStatusPending {
			inDegree[node.Id] = 0
		}
	}

	for _, edge := range w.Edges {
		if !w.isFinalState(w.nodeMap[edge.Source]) && w.nodeMap[edge.Target].Status == NodeStatusPending {
			inDegree[edge.Target]++
		}
	}

	var readyNode []*Node
	for node, degree := range inDegree {
		if degree == 0 {
			w.nodeMap[node].Status = NodeStatusReady
			readyNode = append(readyNode, w.nodeMap[node])
		}
	}
	return readyNode
}

func (w *Workflow) HasCycle() ([]*Node, bool) {
	visited := make(map[*Node]bool)
	recursionStack := make(map[*Node]bool)

	for _, task := range w.Nodes {
		var chain []*Node
		if w.dfs(task, visited, recursionStack, &chain) {
			return chain, true
		}
	}
	return nil, false
}

func (w *Workflow) dfs(task *Node, visited map[*Node]bool, recursionStack map[*Node]bool, chain *[]*Node) bool {
	visited[task] = true
	recursionStack[task] = true
	defer delete(recursionStack, task)

	for _, edge := range w.Edges {
		if w.nodeMap[edge.Source] == task {
			*chain = append(*chain, task)
			if recursionStack[w.nodeMap[edge.Target]] {
				*chain = append(*chain, w.nodeMap[edge.Target])
				return true
			}
			if !visited[w.nodeMap[edge.Target]] {
				if w.dfs(w.nodeMap[edge.Target], visited, recursionStack, chain) {
					return true
				}
			}
		}
	}
	return false
}
