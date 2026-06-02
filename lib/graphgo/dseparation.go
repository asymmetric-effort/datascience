package graphgo

// DSeparation tests whether node sets x and y are d-separated given the
// observed set z in the directed graph g, using the Bayes-Ball algorithm.
//
// Returns true if x and y are d-separated (conditionally independent) given z.
func DSeparation(g *DiGraph, x, y, z map[string]bool) bool {
	// Build the observed set for quick lookup.
	observed := make(map[string]bool, len(z))
	for n := range z {
		observed[n] = true
	}

	// Track visited states: (node, direction). Direction is "from_child" or
	// "from_parent". We encode this as a struct key.
	type visit struct {
		node      string
		fromChild bool // true = arrived from a child, false = arrived from a parent
	}

	visited := make(map[visit]bool)

	// Reachable collects all nodes the ball can reach.
	reachable := make(map[string]bool)

	// BFS / queue-based Bayes-Ball traversal.
	queue := make([]visit, 0, len(x))
	for node := range x {
		// Start from x nodes as if visiting from a child (the ball originates
		// at x, which behaves like arriving from a child for traversal rules).
		v := visit{node: node, fromChild: true}
		queue = append(queue, v)
		visited[v] = true
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		node := curr.node

		reachable[node] = true

		if curr.fromChild {
			// Arrived from a child node.
			if !observed[node] {
				// Not observed: ball passes through to parents and bounces
				// back to other children.
				for p := range g.pred[node] {
					v := visit{node: p, fromChild: true}
					if !visited[v] {
						visited[v] = true
						queue = append(queue, v)
					}
				}
				for c := range g.succ[node] {
					v := visit{node: c, fromChild: false}
					if !visited[v] {
						visited[v] = true
						queue = append(queue, v)
					}
				}
			}
			// If observed: ball is blocked (cannot pass through).
		} else {
			// Arrived from a parent node.
			if !observed[node] {
				// Not observed: ball passes through to children.
				for c := range g.succ[node] {
					v := visit{node: c, fromChild: false}
					if !visited[v] {
						visited[v] = true
						queue = append(queue, v)
					}
				}
			}
			if observed[node] {
				// Observed: v-structure / explaining away. Ball bounces
				// back to parents.
				for p := range g.pred[node] {
					v := visit{node: p, fromChild: true}
					if !visited[v] {
						visited[v] = true
						queue = append(queue, v)
					}
				}
			}
		}
	}

	// If any y node is reachable, x and y are NOT d-separated.
	for node := range y {
		if reachable[node] {
			return false
		}
	}
	return true
}

// MarkovBlanket returns the Markov blanket of a node in a directed graph:
// parents + children + parents of children (co-parents).
func MarkovBlanket(g *DiGraph, node string) map[string]bool {
	blanket := make(map[string]bool)

	// Parents.
	for p := range g.pred[node] {
		blanket[p] = true
	}

	// Children.
	for c := range g.succ[node] {
		blanket[c] = true

		// Parents of children (co-parents).
		for cp := range g.pred[c] {
			if cp != node {
				blanket[cp] = true
			}
		}
	}

	return blanket
}
