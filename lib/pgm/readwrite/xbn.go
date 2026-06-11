package readwrite

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

// XBN XML structures for Microsoft XBN format.
type xbnAnalysisNotebook struct {
	XMLName xml.Name `xml:"ANALYSISNOTEBOOK"`
	BNMODEL xbnModel `xml:"BNMODEL"`
}

type xbnModel struct {
	Name       string        `xml:"NAME,attr"`
	StaticProp xbnStaticProp `xml:"STATICPROPERTIES"`
	DynaProp   xbnDynaProp   `xml:"DYNAMICPROPERTIES"`
}

type xbnStaticProp struct {
	NodeList xbnNodeList `xml:"NODELIST"`
	ArcList  xbnArcList  `xml:"ARCLIST"`
}

type xbnNodeList struct {
	Nodes []xbnNode `xml:"NODE"`
}

type xbnNode struct {
	Name   string     `xml:"NAME,attr"`
	States []xbnState `xml:"STATENAME"`
}

type xbnState struct {
	Value string `xml:",chardata"`
}

type xbnArcList struct {
	Arcs []xbnArc `xml:"ARC"`
}

type xbnArc struct {
	Parent string `xml:"PARENT,attr"`
	Child  string `xml:"CHILD,attr"`
}

type xbnDynaProp struct {
	Formats []xbnFormat `xml:"FORMAT"`
	Dists   []xbnDist   `xml:"DISTRIBS>DIST"`
}

type xbnFormat struct {
	// Placeholder; not used in parsing.
}

type xbnDist struct {
	Type     string        `xml:"TYPE,attr"`
	CondElem []xbnCondElem `xml:"CONDSET>CONDELEM"`
	DPIs     []xbnDPI      `xml:"DPIS>DPI"`
	PrivDist []xbnPrivDist `xml:"PRIVATE>DPIS>DPI"`
}

type xbnCondElem struct {
	Name string `xml:"NAME,attr"`
}

type xbnDPI struct {
	Indexes string `xml:"INDEXES,attr,omitempty"`
	Values  string `xml:",chardata"`
}

type xbnPrivDist struct {
	Indexes string `xml:"INDEXES,attr,omitempty"`
	Values  string `xml:",chardata"`
}

// ReadXBN parses a Microsoft XBN format file and returns a BayesianNetwork.
// Supports full NODELIST with STATENAME elements, ARCLIST with parent/child
// arcs, and DISTRIBS with DIST elements containing CONDSET/CONDELEM for parent
// references and DPIS/DPI with INDEXES attributes for conditional probability
// distributions.
func ReadXBN(r io.Reader) (*models.BayesianNetwork, error) {
	return ReadXBNWithLimit(r, MaxInputSize)
}

// ReadXBNWithLimit is like ReadXBN but accepts a custom maximum input size
// in bytes. Use this for models larger than MaxInputSize (1 MB).
func ReadXBNWithLimit(r io.Reader, maxBytes int) (*models.BayesianNetwork, error) {
	bn := models.NewBayesianNetwork()
	if err := readXBNWith(r, &realBuilder{bn: bn}, bn, maxBytes); err != nil {
		return nil, err
	}
	return bn, nil
}

// readXBNWith is the testable implementation of ReadXBN. Accepts a bnBuilder
// interface for mock injection.
func readXBNWith(r io.Reader, builder bnBuilder, bn *models.BayesianNetwork, maxBytes int) error {
	data, err := readLimitedN(r, maxBytes)
	if err != nil {
		return fmt.Errorf("readwrite: error reading XBN: %w", err)
	}

	var notebook xbnAnalysisNotebook
	if err := xml.Unmarshal(data, &notebook); err != nil {
		return fmt.Errorf("readwrite: error parsing XBN: %w", err)
	}

	type varInfo struct {
		card   int
		states []string
	}
	varMap := make(map[string]*varInfo)

	// Add nodes.
	var nodeOrder []string
	for _, node := range notebook.BNMODEL.StaticProp.NodeList.Nodes {
		name := node.Name
		states := make([]string, len(node.States))
		for i, s := range node.States {
			states[i] = strings.TrimSpace(s.Value)
		}
		if len(states) == 0 {
			states = []string{"s0", "s1"} // default binary
		}
		if err := builder.AddNode(name); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
		if err := builder.SetStates(name, states); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
		varMap[name] = &varInfo{card: len(states), states: states}
		nodeOrder = append(nodeOrder, name)
	}

	// Add arcs.
	for _, arc := range notebook.BNMODEL.StaticProp.ArcList.Arcs {
		if err := builder.AddEdge(arc.Parent, arc.Child); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("readwrite: %w", err)
			}
		}
	}

	// Parse distributions.
	for i, dist := range notebook.BNMODEL.DynaProp.Dists {
		if i >= len(nodeOrder) {
			break
		}
		child := nodeOrder[i]
		childInfo := varMap[child]

		var parents []string
		var evidenceCard []int
		for _, ce := range dist.CondElem {
			p := ce.Name
			if p == child {
				continue
			}
			pi := varMap[p]
			if pi == nil {
				continue
			}
			parents = append(parents, p)
			evidenceCard = append(evidenceCard, pi.card)
		}

		numParentConfigs, err := safeParentConfigs(evidenceCard)
		if err != nil {
			return fmt.Errorf("readwrite: XBN distribution for %q: %w", child, err)
		}

		// Initialize values array: values[childState][parentConfig].
		values := make([][]float64, childInfo.card)
		for cs := 0; cs < childInfo.card; cs++ {
			values[cs] = make([]float64, numParentConfigs)
		}

		if len(parents) == 0 {
			// Unconditional: single DPI row with all child state probs.
			var allVals []float64
			for _, dpi := range dist.DPIs {
				vals, err := xmlbifParseFloats(dpi.Values)
				if err != nil {
					return fmt.Errorf("readwrite: error parsing XBN dist for %q: %w", child, err)
				}
				allVals = append(allVals, vals...)
			}

			if len(allVals) != childInfo.card {
				// Fallback to uniform.
				prob := 1.0 / float64(childInfo.card)
				for cs := 0; cs < childInfo.card; cs++ {
					values[cs][0] = prob
				}
			} else {
				for cs := 0; cs < childInfo.card; cs++ {
					values[cs][0] = allVals[cs]
				}
			}
		} else {
			// Conditional: DPI rows with INDEXES attribute specifying parent config.
			// Check if INDEXES attributes are present.
			hasIndexes := false
			for _, dpi := range dist.DPIs {
				if strings.TrimSpace(dpi.Indexes) != "" {
					hasIndexes = true
					break
				}
			}

			if hasIndexes {
				// Parse DPIs using INDEXES to determine parent config.
				for _, dpi := range dist.DPIs {
					vals, err := xmlbifParseFloats(dpi.Values)
					if err != nil {
						return fmt.Errorf("readwrite: error parsing XBN dist for %q: %w", child, err)
					}
					if len(vals) != childInfo.card {
						continue
					}

					// Parse INDEXES: space-separated integers representing parent
					// state indices. Compute parent config from these.
					idxFields := strings.Fields(strings.TrimSpace(dpi.Indexes))
					if len(idxFields) != len(parents) {
						continue
					}

					pc := 0
					stride := 1
					valid := true
					for pi := len(parents) - 1; pi >= 0; pi-- {
						idx, err := strconv.Atoi(idxFields[pi])
						if err != nil {
							valid = false
							break
						}
						pc += idx * stride
						stride *= evidenceCard[pi]
					}
					if !valid {
						continue
					}

					if pc < numParentConfigs {
						for cs := 0; cs < childInfo.card; cs++ {
							values[cs][pc] = vals[cs]
						}
					}
				}
			} else {
				// No INDEXES: DPI rows are in order of parent configs.
				var allVals []float64
				for _, dpi := range dist.DPIs {
					vals, err := xmlbifParseFloats(dpi.Values)
					if err != nil {
						return fmt.Errorf("readwrite: error parsing XBN dist for %q: %w", child, err)
					}
					allVals = append(allVals, vals...)
				}

				expectedLen := childInfo.card * numParentConfigs
				if len(allVals) != expectedLen {
					// Fallback to uniform.
					prob := 1.0 / float64(childInfo.card)
					for cs := 0; cs < childInfo.card; cs++ {
						for pc := 0; pc < numParentConfigs; pc++ {
							values[cs][pc] = prob
						}
					}
				} else {
					idx := 0
					for pc := 0; pc < numParentConfigs; pc++ {
						for cs := 0; cs < childInfo.card; cs++ {
							values[cs][pc] = allVals[idx]
							idx++
						}
					}
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

	// For nodes without dists, create uniform CPDs.
	for _, name := range nodeOrder {
		if bn.GetCPD(name) != nil {
			continue
		}
		info := varMap[name]
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

	return nil
}

// WriteXBN serializes a BayesianNetwork to Microsoft XBN format with full
// NODELIST (including STATENAME elements), ARCLIST, and DISTRIBS with
// CONDSET/CONDELEM and DPI elements with INDEXES attributes.
func WriteXBN(w io.Writer, bn *models.BayesianNetwork) error {
	nodes := bn.Nodes()

	// Build node list.
	var xbnNodes []xbnNode
	for _, node := range nodes {
		states := bn.GetStates(node)
		if len(states) == 0 {
			return fmt.Errorf("readwrite: variable %q has no state names", node)
		}
		xStates := make([]xbnState, len(states))
		for i, s := range states {
			xStates[i] = xbnState{Value: s}
		}
		xbnNodes = append(xbnNodes, xbnNode{Name: node, States: xStates})
	}

	// Build arc list.
	edges := bn.Edges()
	arcs := make([]xbnArc, len(edges))
	for i, e := range edges {
		arcs[i] = xbnArc{Parent: e[0], Child: e[1]}
	}

	// Build distributions.
	var dists []xbnDist
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

		var condElems []xbnCondElem
		for _, ev := range evidence {
			condElems = append(condElems, xbnCondElem{Name: ev})
		}

		var dpis []xbnDPI
		for pc := 0; pc < numParentConfigs; pc++ {
			var parts []string
			for cs := 0; cs < childCard; cs++ {
				parts = append(parts, formatFloat(data[cs*numParentConfigs+pc]))
			}

			dpi := xbnDPI{Values: strings.Join(parts, " ")}

			// Add INDEXES for conditional distributions.
			if len(evidence) > 0 {
				indexParts := make([]string, len(evidence))
				remainder := pc
				for pi := len(evidence) - 1; pi >= 0; pi-- {
					indexParts[pi] = strconv.Itoa(remainder % evidenceCard[pi])
					remainder /= evidenceCard[pi]
				}
				dpi.Indexes = strings.Join(indexParts, " ")
			}

			dpis = append(dpis, dpi)
		}

		dists = append(dists, xbnDist{
			Type:     "discrete",
			CondElem: condElems,
			DPIs:     dpis,
		})
	}

	notebook := xbnAnalysisNotebook{
		BNMODEL: xbnModel{
			Name: "unknown",
			StaticProp: xbnStaticProp{
				NodeList: xbnNodeList{Nodes: xbnNodes},
				ArcList:  xbnArcList{Arcs: arcs},
			},
			DynaProp: xbnDynaProp{
				Dists: dists,
			},
		},
	}

	if _, err := fmt.Fprint(w, xml.Header); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(notebook); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}
	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	return nil
}
