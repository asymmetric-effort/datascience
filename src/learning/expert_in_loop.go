package learning

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ExpertInLoop implements LLM-in-the-loop causal discovery. It combines
// statistical conditional independence testing (for skeleton discovery) with
// LLM-based causal reasoning (for edge orientation). When the PC algorithm
// would orient a v-structure using only the separating set, ExpertInLoop also
// queries the LLM for its opinion on causal direction and combines the two
// signals.
type ExpertInLoop struct {
	data         *tabgo.DataFrame
	llmClient    LLMClient
	ciTest       CITestFunc
	significance float64
}

// NewExpertInLoop creates a new ExpertInLoop causal discovery instance.
//
// Parameters:
//   - data: the observational dataset with columns as variables
//   - llmClient: an LLM client for querying causal direction opinions
//   - ciTest: a conditional independence test function
//   - significance: the significance level (alpha) for independence tests
func NewExpertInLoop(data *tabgo.DataFrame, llmClient LLMClient, ciTest CITestFunc, significance float64) *ExpertInLoop {
	return &ExpertInLoop{
		data:         data,
		llmClient:    llmClient,
		ciTest:       ciTest,
		significance: significance,
	}
}

// Estimate runs LLM-in-the-loop causal discovery and returns a PDAG.
//
// The algorithm proceeds in four phases:
//  1. Skeleton discovery via conditional independence tests (same as PC)
//  2. V-structure orientation combining statistical evidence with LLM opinion
//  3. Meek rule application to propagate orientations
//  4. Return the resulting PDAG
func (e *ExpertInLoop) Estimate() (*graphgo.PDAG, error) {
	variables := e.data.Columns()
	if len(variables) < 2 {
		return nil, fmt.Errorf("learning: ExpertInLoop requires at least 2 variables, got %d", len(variables))
	}

	// Phase 1: Skeleton discovery (identical to PC algorithm).
	pdag := graphgo.NewPDAG()
	for _, v := range variables {
		pdag.AddNode(v)
	}
	for i := 0; i < len(variables); i++ {
		for j := i + 1; j < len(variables); j++ {
			pdag.AddUndirectedEdge(variables[i], variables[j])
		}
	}

	sepSets := make(map[[2]string][]string)

	for d := 0; ; d++ {
		maxDegree := 0
		for _, node := range variables {
			deg := e.adjacencyCount(pdag, node)
			if deg > maxDegree {
				maxDegree = deg
			}
		}
		if d > maxDegree {
			break
		}

		edges := pdag.UndirectedEdges()
		for _, edge := range edges {
			x, y := edge[0], edge[1]
			if !pdag.HasUndirectedEdge(x, y) {
				continue
			}

			adjX := e.undirectedNeighbors(pdag, x)
			adjXMinusY := removeFromSlice(adjX, y)
			if found, subset := e.findSepSet(x, y, adjXMinusY, d); found {
				pdag.RemoveUndirectedEdge(x, y)
				sepSets[sepSetKey(x, y)] = subset
				continue
			}

			adjY := e.undirectedNeighbors(pdag, y)
			adjYMinusX := removeFromSlice(adjY, x)
			if found, subset := e.findSepSet(x, y, adjYMinusX, d); found {
				pdag.RemoveUndirectedEdge(x, y)
				sepSets[sepSetKey(x, y)] = subset
				continue
			}
		}
	}

	// Phase 2: V-structure orientation with LLM consultation.
	// For each unshielded triple x - z - y (where x and y are not adjacent),
	// combine the statistical signal (z not in sepSet) with the LLM opinion.
	promptTemplate := CausalPromptTemplate{}

	for _, z := range variables {
		neighborsZ := e.undirectedNeighbors(pdag, z)
		if len(neighborsZ) < 2 {
			continue
		}
		for i := 0; i < len(neighborsZ); i++ {
			for j := i + 1; j < len(neighborsZ); j++ {
				x, y := neighborsZ[i], neighborsZ[j]
				if pdag.Adjacent(x, y) {
					continue
				}

				key := sepSetKey(x, y)
				ss, exists := sepSets[key]
				if !exists {
					continue
				}

				// Statistical signal: z NOT in sepSet suggests v-structure.
				statisticalVStructure := !containsString(ss, z)

				// Query LLM for causal direction opinion.
				llmVStructure := e.queryLLMForVStructure(promptTemplate, x, y, z)

				// Combine signals: orient as v-structure if either signal
				// supports it and neither contradicts. The statistical test
				// is the primary signal; the LLM can reinforce or break ties.
				orient := e.combineSignals(statisticalVStructure, llmVStructure)

				if orient {
					if pdag.HasUndirectedEdge(x, z) {
						pdag.RemoveUndirectedEdge(x, z)
						pdag.AddDirectedEdge(x, z)
					}
					if pdag.HasUndirectedEdge(y, z) {
						pdag.RemoveUndirectedEdge(y, z)
						pdag.AddDirectedEdge(y, z)
					}
				}
			}
		}
	}

	// Phase 3: Meek rules.
	graphgo.ApplyMeekRules(pdag)

	return pdag, nil
}

// llmOpinion represents the LLM's opinion on a v-structure.
type llmOpinion int

const (
	llmSupports  llmOpinion = iota // LLM says YES to v-structure
	llmOpposes                     // LLM says NO to v-structure
	llmUncertain                   // LLM says UNKNOWN or is unavailable
)

// queryLLMForVStructure queries the LLM about whether x -> z <- y is a
// v-structure. It asks two questions: "does x cause z?" and "does y cause z?"
// If both are YES, the LLM supports the v-structure. If either is NO, it
// opposes. Otherwise, it is uncertain.
func (e *ExpertInLoop) queryLLMForVStructure(tmpl CausalPromptTemplate, x, y, z string) llmOpinion {
	if e.llmClient == nil {
		return llmUncertain
	}

	// Build context from the variable names involved.
	context := fmt.Sprintf("variables %s, %s, and %s in a causal graph where %s and %s are not directly connected but both connected to %s",
		x, y, z, x, y, z)

	// Query: does x cause z?
	dirXZ := e.queryOneCausalDirection(tmpl, x, z, context)
	// Query: does y cause z?
	dirYZ := e.queryOneCausalDirection(tmpl, y, z, context)

	if dirXZ == "YES" && dirYZ == "YES" {
		return llmSupports
	}
	if dirXZ == "NO" || dirYZ == "NO" {
		return llmOpposes
	}
	return llmUncertain
}

// queryOneCausalDirection queries the LLM for a single causal direction.
// Returns "YES", "NO", or "UNKNOWN". On any error, returns "UNKNOWN".
func (e *ExpertInLoop) queryOneCausalDirection(tmpl CausalPromptTemplate, from, to, context string) string {
	prompt := tmpl.FormatCausalDirectionQuery(from, to, context)
	response, err := e.llmClient.Complete(prompt, WithTemperature(0.0))
	if err != nil {
		return "UNKNOWN"
	}

	direction, _, err := tmpl.ParseCausalResponse(response)
	if err != nil {
		return "UNKNOWN"
	}
	return direction
}

// combineSignals combines the statistical v-structure decision with the LLM
// opinion. The rules are:
//   - If statistical evidence says v-structure AND LLM supports or is uncertain: orient
//   - If statistical evidence says v-structure BUT LLM opposes: still orient
//     (statistical evidence is primary, but this could be made configurable)
//   - If statistical evidence says NOT v-structure AND LLM supports: orient
//     (LLM can override when stats are ambiguous)
//   - If statistical evidence says NOT v-structure AND LLM opposes or is uncertain: don't orient
func (e *ExpertInLoop) combineSignals(statisticalVStructure bool, llm llmOpinion) bool {
	if statisticalVStructure {
		// Statistical evidence is primary; orient regardless of LLM.
		return true
	}
	// No statistical evidence for v-structure.
	// LLM can promote a non-v-structure to v-structure if it strongly supports.
	if llm == llmSupports {
		return true
	}
	return false
}

// EstimateBN runs ExpertInLoop via Estimate(), converts the resulting PDAG to
// a DAG, and returns a BayesianNetwork (without CPDs — only structure).
func (e *ExpertInLoop) EstimateBN() (*models.BayesianNetwork, error) {
	pdag, err := e.Estimate()
	if err != nil {
		return nil, err
	}

	dag, err := pdagToDAG(pdag)
	if err != nil {
		return nil, fmt.Errorf("learning: ExpertInLoop failed to convert PDAG to DAG: %w", err)
	}
	return dag, nil
}

// adjacencyCount returns the number of undirected neighbors of node in the PDAG.
func (e *ExpertInLoop) adjacencyCount(pdag *graphgo.PDAG, node string) int {
	count := 0
	for _, n := range pdag.Nodes() {
		if n != node && pdag.HasUndirectedEdge(node, n) {
			count++
		}
	}
	return count
}

// undirectedNeighbors returns sorted undirected neighbors of node in the PDAG.
func (e *ExpertInLoop) undirectedNeighbors(pdag *graphgo.PDAG, node string) []string {
	var neighbors []string
	for _, n := range pdag.Nodes() {
		if n != node && pdag.HasUndirectedEdge(node, n) {
			neighbors = append(neighbors, n)
		}
	}
	sort.Strings(neighbors)
	return neighbors
}

// findSepSet searches for a subset of candidates of size d that makes x and y
// conditionally independent.
func (e *ExpertInLoop) findSepSet(x, y string, candidates []string, d int) (bool, []string) {
	if d > len(candidates) {
		return false, nil
	}
	if d == 0 {
		_, _, indep := e.ciTest(x, y, nil, e.data, e.significance)
		if indep {
			return true, nil
		}
		return false, nil
	}

	subsets := combinations(candidates, d)
	for _, subset := range subsets {
		_, _, indep := e.ciTest(x, y, subset, e.data, e.significance)
		if indep {
			return true, subset
		}
	}
	return false, nil
}

// buildVariableContext constructs a contextual description of the variables
// involved for LLM prompting. This is kept simple; a production system would
// incorporate domain metadata.
func buildVariableContext(vars []string) string {
	return "observational data with variables: " + strings.Join(vars, ", ")
}
