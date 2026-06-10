package readwrite

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

// jsonNetwork is the JSON representation of a Bayesian network.
type jsonNetwork struct {
	Name   string              `json:"name"`
	Nodes  []string            `json:"nodes"`
	Edges  [][2]string         `json:"edges"`
	States map[string][]string `json:"states,omitempty"`
	CPDs   map[string]*jsonCPD `json:"cpds,omitempty"`
}

// jsonCPD is the JSON representation of a TabularCPD.
type jsonCPD struct {
	VariableCard int         `json:"variable_card"`
	Values       [][]float64 `json:"values"`
	Evidence     []string    `json:"evidence,omitempty"`
	EvidenceCard []int       `json:"evidence_card,omitempty"`
}

// ReadJSON parses a JSON file and returns a fully populated BayesianNetwork,
// including nodes, edges, states, and CPDs.
func ReadJSON(r io.Reader) (*models.BayesianNetwork, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("readwrite: error reading JSON: %w", err)
	}

	var jn jsonNetwork
	if err := json.Unmarshal(data, &jn); err != nil {
		return nil, fmt.Errorf("readwrite: error parsing JSON: %w", err)
	}

	bn, err := jsonBuildStructure(&jn)
	if err != nil {
		return nil, err
	}

	if err := jsonAddCPDs(&jn, &realBuilder{bn: bn}); err != nil {
		return nil, err
	}

	return bn, nil
}

// jsonAddCPDs is the testable implementation for adding CPDs during JSON read.
// Accepts a bnBuilder interface for mock injection.
func jsonAddCPDs(jn *jsonNetwork, builder bnBuilder) error {
	for variable, jcpd := range jn.CPDs {
		cpd, err := factors.NewTabularCPD(variable, jcpd.VariableCard, jcpd.Values,
			jcpd.Evidence, jcpd.EvidenceCard)
		if err != nil {
			return fmt.Errorf("readwrite: failed to create CPD for %q: %w", variable, err)
		}
		if err := builder.AddCPD(cpd); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
	}
	return nil
}

// ReadJSONStructure parses a JSON file and returns a BayesianNetwork with
// structure only (nodes, edges, states). CPDs in the JSON are ignored.
func ReadJSONStructure(r io.Reader) (*models.BayesianNetwork, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("readwrite: error reading JSON: %w", err)
	}

	var jn jsonNetwork
	if err := json.Unmarshal(data, &jn); err != nil {
		return nil, fmt.Errorf("readwrite: error parsing JSON: %w", err)
	}

	return jsonBuildStructure(&jn)
}

// WriteJSON serializes a BayesianNetwork to JSON format.
func WriteJSON(w io.Writer, bn *models.BayesianNetwork) error {
	nodes := bn.Nodes()
	edges := bn.Edges()

	states := make(map[string][]string)
	for _, node := range nodes {
		s := bn.GetStates(node)
		if len(s) > 0 {
			states[node] = s
		}
	}

	cpds := make(map[string]*jsonCPD)
	for _, node := range nodes {
		cpd := bn.GetCPD(node)
		if cpd == nil {
			return fmt.Errorf("readwrite: variable %q has no CPD", node)
		}

		evidence := cpd.Evidence()
		evidenceCard := cpd.EvidenceCard()
		childCard := cpd.VariableCard()

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		data := cpd.ToFactor().Values().Data()

		// Reconstruct values[childState][parentConfig].
		values := make([][]float64, childCard)
		for cs := 0; cs < childCard; cs++ {
			values[cs] = make([]float64, numParentConfigs)
			for pc := 0; pc < numParentConfigs; pc++ {
				values[cs][pc] = data[cs*numParentConfigs+pc]
			}
		}

		jcpd := &jsonCPD{
			VariableCard: childCard,
			Values:       values,
		}
		if len(evidence) > 0 {
			jcpd.Evidence = evidence
			jcpd.EvidenceCard = evidenceCard
		}
		cpds[node] = jcpd
	}

	jn := jsonNetwork{
		Name:   "network",
		Nodes:  nodes,
		Edges:  edges,
		States: states,
		CPDs:   cpds,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jn); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	return nil
}

// jsonBuildStructure creates a BayesianNetwork from parsed JSON (structure only).
func jsonBuildStructure(jn *jsonNetwork) (*models.BayesianNetwork, error) {
	bn := models.NewBayesianNetwork()

	for _, node := range jn.Nodes {
		node = strings.TrimSpace(node)
		if err := bn.AddNode(node); err != nil {
			return nil, fmt.Errorf("readwrite: %w", err)
		}
	}

	for _, edge := range jn.Edges {
		from := strings.TrimSpace(edge[0])
		to := strings.TrimSpace(edge[1])
		if err := bn.AddEdge(from, to); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return nil, fmt.Errorf("readwrite: %w", err)
			}
		}
	}

	for variable, stateNames := range jn.States {
		if err := bn.SetStates(variable, stateNames); err != nil {
			return nil, fmt.Errorf("readwrite: %w", err)
		}
	}

	return bn, nil
}
