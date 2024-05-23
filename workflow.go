package workflow

import "github.com/sirupsen/logrus"

func (w *Workflow) AddNode(node ...*Node) {
	for _, n := range node {
		n.Status = NodeStatusPending
	}
	w.Nodes = append(w.Nodes, node...)
}

func (w *Workflow) AddEdge(source, target *Node) {
	w.Edges = append(w.Edges, &Edge{Source: source, Target: target})
}

func (w *Workflow) IsNil() bool {
	return w == nil || len(w.Nodes) == 0
}

func (w *Workflow) ReplaceNodeToSubWorkflow(source *Node, target *Workflow) {
	// 替换节点
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
	find := make(map[string]*Node)
	for _, t := range target.Nodes {
		temps[t.Id] = 0
		tempt[t.Id] = 0
		find[t.Id] = t
	}

	for _, e := range target.Edges {
		temps[e.Source.Id]++
		tempt[e.Target.Id]++
	}

	for index := len(w.Edges) - 1; index >= 0; index-- {
		edge := w.Edges[index]
		if edge.Source.Id == source.Id {
			var ts []*Node
			for t, degree := range temps {
				if degree == 0 {
					ts = append(ts, find[t])
				}
			}
			for _, t := range ts {
				w.Edges = append(w.Edges, &Edge{Source: t, Target: edge.Target})
			}

			// 删除边
			w.Edges = append(w.Edges[:index], w.Edges[index+1:]...)
		}
		if edge.Target.Id == source.Id {
			var ts []*Node
			for t, degree := range tempt {
				if degree == 0 {
					ts = append(ts, find[t])
				}
			}

			for _, t := range ts {
				w.Edges = append(w.Edges, &Edge{Source: edge.Source, Target: t})
			}

			// 删除边
			w.Edges = append(w.Edges[:index], w.Edges[index+1:]...)
		}
	}
}

func (w *Workflow) GetWorkflow() *Workflow {
	for _, edge := range w.Edges {
		edge.Status = edge.Target.Status
	}
	return w
}

func (w *Workflow) isFinalState(node *Node) bool {
	return node.Status == NodeStatusSuccess || node.Status == NodeStatusFailure || node.Status == NodeStatusAborted
}

func (w *Workflow) GetReadyNodes() []*Node {
	inDegree := make(map[string]int)
	nodeMap := make(map[string]*Node)
	for _, node := range w.Nodes {
		if node.Status == NodeStatusPending {
			inDegree[node.Id] = 0
			nodeMap[node.Id] = node
		}
	}

	for _, edge := range w.Edges {
		logrus.Infof("edge: %+v", edge)
		if !w.isFinalState(edge.Source) && edge.Target.Status == NodeStatusPending {
			inDegree[edge.Target.Id]++
		}
	}

	var readyNode []*Node
	for node, degree := range inDegree {
		if degree == 0 {
			nodeMap[node].Status = NodeStatusReady
			readyNode = append(readyNode, nodeMap[node])
		}
	}
	logrus.Infof("ready node: %+v", readyNode)
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
		if edge.Source == task {
			*chain = append(*chain, task)
			if recursionStack[edge.Target] {
				*chain = append(*chain, edge.Target)
				return true
			}
			if !visited[edge.Target] {
				if w.dfs(edge.Target, visited, recursionStack, chain) {
					return true
				}
			}
		}
	}
	return false
}
