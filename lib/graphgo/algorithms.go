package graphgo

import "fmt"

// IsDAG returns true if the directed graph is a directed acyclic graph.
func IsDAG(g *DiGraph) bool {
	_, err := TopologicalSort(g)
	return err == nil
}

// TopologicalSort returns a topological ordering of the DAG using Kahn's algorithm.
// Returns an error if the graph contains a cycle.
func TopologicalSort(g *DiGraph) ([]string, error) {
	inDegree := make(map[string]int, len(g.succ))
	for n := range g.succ {
		inDegree[n] = len(g.pred[n])
	}

	// Seed queue with zero in-degree nodes.
	queue := make([]string, 0)
	for n, d := range inDegree {
		if d == 0 {
			queue = append(queue, n)
		}
	}

	var order []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for succ := range g.succ[node] {
			inDegree[succ]--
			if inDegree[succ] == 0 {
				queue = append(queue, succ)
			}
		}
	}

	if len(order) != len(g.succ) {
		return nil, fmt.Errorf("graphgo: graph contains a cycle")
	}
	return order, nil
}

// Ancestors returns the set of all nodes that have a directed path to the given node.
func Ancestors(g *DiGraph, node string) map[string]bool {
	visited := make(map[string]bool)
	var dfs func(string)
	dfs = func(n string) {
		for p := range g.pred[n] {
			if !visited[p] {
				visited[p] = true
				dfs(p)
			}
		}
	}
	dfs(node)
	return visited
}

// Descendants returns the set of all nodes reachable from the given node via directed edges.
func Descendants(g *DiGraph, node string) map[string]bool {
	visited := make(map[string]bool)
	var dfs func(string)
	dfs = func(n string) {
		for s := range g.succ[n] {
			if !visited[s] {
				visited[s] = true
				dfs(s)
			}
		}
	}
	dfs(node)
	return visited
}
