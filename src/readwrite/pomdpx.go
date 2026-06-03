package readwrite

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// PomdpX XML structures with full CondProb/Entry support for BN use.
type pomdpxDoc struct {
	XMLName         xml.Name          `xml:"pomdpx"`
	Version         string            `xml:"version,attr"`
	ID              string            `xml:"id,attr,omitempty"`
	Description     string            `xml:"Description,omitempty"`
	Variables       pomdpxVarBlock    `xml:"Variable"`
	InitBelief      *pomdpxInit       `xml:"InitialStateBelief,omitempty"`
	StateTransition *pomdpxTransBlock `xml:"StateTransitionFunction,omitempty"`
}

type pomdpxVarBlock struct {
	StateVars []pomdpxStateVar `xml:"StateVar"`
}

type pomdpxStateVar struct {
	VarName    string   `xml:"vnamePrev,attr"`
	NumValues  int      `xml:"numValues,attr,omitempty"`
	ValueNames []string `xml:"ValueEnum,omitempty"`
}

type pomdpxInit struct {
	CondProbs []pomdpxCondProb `xml:"CondProb"`
}

type pomdpxTransBlock struct {
	CondProbs []pomdpxCondProb `xml:"CondProb"`
}

type pomdpxCondProb struct {
	Name    string         `xml:"name,attr,omitempty"`
	Var     []pomdpxVar    `xml:"Var"`
	Parents []pomdpxParent `xml:"Parent"`
	Params  []pomdpxParam  `xml:"Parameter"`
}

type pomdpxVar struct {
	Name string `xml:",chardata"`
}

type pomdpxParent struct {
	Names string `xml:",chardata"`
}

type pomdpxParam struct {
	Type    string        `xml:"type,attr,omitempty"`
	Entries []pomdpxEntry `xml:"Entry"`
}

type pomdpxEntry struct {
	Instance  string `xml:"Instance,omitempty"`
	ProbTable string `xml:"ProbTable"`
}

// ReadPomdpX parses a PomdpX format file and returns a BayesianNetwork.
// Supports full Variable definitions with ValueEnum, InitialStateBelief for
// unconditional CPDs, and StateTransitionFunction for conditional CPDs with
// parent references via Entry instance/probability pairs.
func ReadPomdpX(r io.Reader) (*models.BayesianNetwork, error) {
	bn := models.NewBayesianNetwork()
	if err := readPomdpXWith(r, &realBuilder{bn: bn}, bn); err != nil {
		return nil, err
	}
	return bn, nil
}

// readPomdpXWith is the testable implementation of ReadPomdpX. Accepts a
// bnBuilder interface for mock injection.
func readPomdpXWith(r io.Reader, builder bnBuilder, bn *models.BayesianNetwork) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("readwrite: error reading PomdpX: %w", err)
	}

	var doc pomdpxDoc
	if err := xml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("readwrite: error parsing PomdpX: %w", err)
	}

	type varInfo struct {
		card   int
		states []string
	}
	varMap := make(map[string]*varInfo)

	// Add state variables.
	for _, sv := range doc.Variables.StateVars {
		name := sv.VarName
		if name == "" {
			continue
		}

		var states []string
		if len(sv.ValueNames) > 0 {
			// ValueEnum is space-separated in the XML.
			for _, ve := range sv.ValueNames {
				for _, s := range strings.Fields(ve) {
					states = append(states, s)
				}
			}
		}

		numVals := sv.NumValues
		if numVals <= 0 {
			numVals = len(states)
		}
		if len(states) == 0 {
			// Generate default state names.
			states = make([]string, numVals)
			for i := 0; i < numVals; i++ {
				states[i] = fmt.Sprintf("s%d", i)
			}
		}

		if err := builder.AddNode(name); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
		if err := builder.SetStates(name, states); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
		varMap[name] = &varInfo{card: len(states), states: states}
	}

	// parseCondProb processes a CondProb element, adds edges and creates a CPD.
	parseCondProb := func(cp pomdpxCondProb) error {
		if len(cp.Var) == 0 {
			return nil
		}
		child := strings.TrimSpace(cp.Var[0].Name)
		childInfo := varMap[child]
		if childInfo == nil {
			return nil
		}

		// Determine parents from the Parent element.
		var parents []string
		var evidenceCard []int
		for _, p := range cp.Parents {
			for _, pName := range strings.Fields(p.Names) {
				pName = strings.TrimSpace(pName)
				if pName == "" || pName == child {
					continue
				}
				pi := varMap[pName]
				if pi == nil {
					return fmt.Errorf("readwrite: PomdpX CondProb references unknown parent %q", pName)
				}
				parents = append(parents, pName)
				evidenceCard = append(evidenceCard, pi.card)
			}
		}

		// Add edges.
		for _, p := range parents {
			if err := builder.AddEdge(p, child); err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					return fmt.Errorf("readwrite: %w", err)
				}
			}
		}

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		if len(parents) == 0 {
			// Unconditional: collect all probability values from entries.
			var allVals []float64
			for _, param := range cp.Params {
				for _, entry := range param.Entries {
					vals, err := xmlbifParseFloats(entry.ProbTable)
					if err != nil {
						return fmt.Errorf("readwrite: error parsing PomdpX probs for %q: %w", child, err)
					}
					allVals = append(allVals, vals...)
				}
			}

			if len(allVals) == childInfo.card {
				values := make([][]float64, childInfo.card)
				for cs := 0; cs < childInfo.card; cs++ {
					values[cs] = []float64{allVals[cs]}
				}

				cpd, err := factors.NewTabularCPD(child, childInfo.card, values, nil, nil)
				if err != nil {
					return fmt.Errorf("readwrite: failed to create CPD for %q: %w", child, err)
				}
				if err := builder.AddCPD(cpd); err != nil {
					return fmt.Errorf("readwrite: %w", err)
				}
			}
		} else {
			// Conditional: parse entries with Instance/ProbTable pairs.
			// Build state-to-index maps for parents.
			parentStateIdx := make([]map[string]int, len(parents))
			for i, p := range parents {
				pi := varMap[p]
				parentStateIdx[i] = make(map[string]int, pi.card)
				for j, s := range pi.states {
					parentStateIdx[i][s] = j
				}
			}

			// Initialize values array: values[childState][parentConfig].
			values := make([][]float64, childInfo.card)
			for cs := 0; cs < childInfo.card; cs++ {
				values[cs] = make([]float64, numParentConfigs)
			}

			// Check if entries use Instance tags or are flat tables.
			hasInstances := false
			for _, param := range cp.Params {
				for _, entry := range param.Entries {
					if strings.TrimSpace(entry.Instance) != "" {
						hasInstances = true
						break
					}
				}
				if hasInstances {
					break
				}
			}

			if hasInstances {
				// Parse entries with instance/probability pairs.
				for _, param := range cp.Params {
					for _, entry := range param.Entries {
						instance := strings.TrimSpace(entry.Instance)
						if instance == "" {
							continue
						}
						instParts := strings.Fields(instance)

						vals, err := xmlbifParseFloats(entry.ProbTable)
						if err != nil {
							return fmt.Errorf("readwrite: error parsing PomdpX probs for %q: %w", child, err)
						}

						if len(vals) != childInfo.card {
							return fmt.Errorf("readwrite: PomdpX entry for %q has %d values, expected %d",
								child, len(vals), childInfo.card)
						}

						// Instance parts are parent state names. Compute the parent config index.
						if len(instParts) != len(parents) {
							return fmt.Errorf("readwrite: PomdpX entry for %q has %d instance parts, expected %d parents",
								child, len(instParts), len(parents))
						}

						pc := 0
						stride := 1
						for pi := len(parents) - 1; pi >= 0; pi-- {
							stateIdx, ok := parentStateIdx[pi][instParts[pi]]
							if !ok {
								return fmt.Errorf("readwrite: PomdpX unknown parent state %q for parent %q",
									instParts[pi], parents[pi])
							}
							pc += stateIdx * stride
							stride *= evidenceCard[pi]
						}

						for cs := 0; cs < childInfo.card; cs++ {
							values[cs][pc] = vals[cs]
						}
					}
				}
			} else {
				// Flat table: all probabilities in order.
				var allVals []float64
				for _, param := range cp.Params {
					for _, entry := range param.Entries {
						vals, err := xmlbifParseFloats(entry.ProbTable)
						if err != nil {
							return fmt.Errorf("readwrite: error parsing PomdpX probs for %q: %w", child, err)
						}
						allVals = append(allVals, vals...)
					}
				}

				expectedLen := childInfo.card * numParentConfigs
				if len(allVals) != expectedLen {
					return fmt.Errorf("readwrite: PomdpX table for %q has %d values, expected %d",
						child, len(allVals), expectedLen)
				}

				// Flat table ordering: for each parent config, list child state probs.
				idx := 0
				for pc := 0; pc < numParentConfigs; pc++ {
					for cs := 0; cs < childInfo.card; cs++ {
						values[cs][pc] = allVals[idx]
						idx++
					}
				}
			}

			cpd, err := factors.NewTabularCPD(child, childInfo.card, values, parents, evidenceCard)
			if err != nil {
				return fmt.Errorf("readwrite: failed to create CPD for %q: %w", child, err)
			}
			if err := builder.AddCPD(cpd); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
		}

		return nil
	}

	// Parse initial belief (unconditional CPDs).
	if doc.InitBelief != nil {
		for _, cp := range doc.InitBelief.CondProbs {
			if err := parseCondProb(cp); err != nil {
				return err
			}
		}
	}

	// Parse state transition function (conditional CPDs).
	if doc.StateTransition != nil {
		for _, cp := range doc.StateTransition.CondProbs {
			if err := parseCondProb(cp); err != nil {
				return err
			}
		}
	}

	// For any variable without a CPD, create a uniform distribution.
	for name, info := range varMap {
		if bn.GetCPD(name) == nil {
			prob := 1.0 / float64(info.card)
			values := make([][]float64, info.card)
			for cs := 0; cs < info.card; cs++ {
				values[cs] = []float64{prob}
			}
			cpd, err := factors.NewTabularCPD(name, info.card, values, nil, nil)
			if err != nil {
				return fmt.Errorf("readwrite: failed to create default CPD for %q: %w", name, err)
			}
			if err := builder.AddCPD(cpd); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
		}
	}

	return nil
}

// WritePomdpX serializes a BayesianNetwork to PomdpX format with full
// conditional probability table support. Unconditional distributions are
// written in InitialStateBelief; conditional distributions are written in
// StateTransitionFunction with Entry elements containing Instance tags.
func WritePomdpX(w io.Writer, bn *models.BayesianNetwork) error {
	nodes := bn.Nodes()

	var stateVars []pomdpxStateVar
	for _, node := range nodes {
		states := bn.GetStates(node)
		if len(states) == 0 {
			return fmt.Errorf("readwrite: variable %q has no state names", node)
		}
		stateVars = append(stateVars, pomdpxStateVar{
			VarName:    node,
			NumValues:  len(states),
			ValueNames: []string{strings.Join(states, " ")},
		})
	}

	// Build initial belief for root nodes and state transition for conditional.
	var initCondProbs []pomdpxCondProb
	var transCondProbs []pomdpxCondProb

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

		if len(evidence) == 0 {
			// Unconditional: initial belief.
			var parts []string
			for cs := 0; cs < childCard; cs++ {
				parts = append(parts, formatFloat(data[cs]))
			}

			initCondProbs = append(initCondProbs, pomdpxCondProb{
				Var: []pomdpxVar{{Name: node}},
				Params: []pomdpxParam{{
					Type: "TBL",
					Entries: []pomdpxEntry{{
						ProbTable: strings.Join(parts, " "),
					}},
				}},
			})
		} else {
			// Conditional: write Entry elements with Instance tags for each
			// parent configuration, mapping parent state names to probabilities.
			var entries []pomdpxEntry
			for pc := 0; pc < numParentConfigs; pc++ {
				// Decompose parent config index into parent state names.
				instanceParts := make([]string, len(evidence))
				remainder := pc
				for pi := len(evidence) - 1; pi >= 0; pi-- {
					stateIdx := remainder % evidenceCard[pi]
					remainder /= evidenceCard[pi]
					parentStates := bn.GetStates(evidence[pi])
					if stateIdx < len(parentStates) {
						instanceParts[pi] = parentStates[stateIdx]
					} else {
						instanceParts[pi] = fmt.Sprintf("s%d", stateIdx)
					}
				}

				var parts []string
				for cs := 0; cs < childCard; cs++ {
					parts = append(parts, formatFloat(data[cs*numParentConfigs+pc]))
				}

				entries = append(entries, pomdpxEntry{
					Instance:  strings.Join(instanceParts, " "),
					ProbTable: strings.Join(parts, " "),
				})
			}

			transCondProbs = append(transCondProbs, pomdpxCondProb{
				Var:     []pomdpxVar{{Name: node}},
				Parents: []pomdpxParent{{Names: strings.Join(evidence, " ")}},
				Params: []pomdpxParam{{
					Type:    "TBL",
					Entries: entries,
				}},
			})
		}
	}

	doc := pomdpxDoc{
		Version:   "1.0",
		Variables: pomdpxVarBlock{StateVars: stateVars},
	}

	if len(initCondProbs) > 0 {
		doc.InitBelief = &pomdpxInit{CondProbs: initCondProbs}
	}
	if len(transCondProbs) > 0 {
		doc.StateTransition = &pomdpxTransBlock{CondProbs: transCondProbs}
	}

	if _, err := fmt.Fprint(w, xml.Header); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}
	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	return nil
}
