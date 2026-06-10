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

// datascience-native XML structures.

type pgmgoNetwork struct {
	XMLName xml.Name   `xml:"pgmgo-network"`
	Name    string     `xml:"name,attr"`
	Nodes   pgmgoNodes `xml:"nodes"`
	Edges   pgmgoEdges `xml:"edges"`
	CPDs    pgmgoCPDs  `xml:"cpds"`
}

type pgmgoNodes struct {
	Nodes []pgmgoNode `xml:"node"`
}

type pgmgoNode struct {
	Name   string `xml:"name,attr"`
	States string `xml:"states,attr"`
}

type pgmgoEdges struct {
	Edges []pgmgoEdge `xml:"edge"`
}

type pgmgoEdge struct {
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
}

type pgmgoCPDs struct {
	CPDs []pgmgoCPD `xml:"cpd"`
}

type pgmgoCPD struct {
	Variable     string `xml:"variable,attr"`
	Card         int    `xml:"card,attr"`
	Evidence     string `xml:"evidence,attr,omitempty"`
	EvidenceCard string `xml:"evidence_card,attr,omitempty"`
	Values       string `xml:"values"`
}

// ReadXMLNative parses a datascience-native XML file and returns a BayesianNetwork.
func ReadXMLNative(r io.Reader) (*models.BayesianNetwork, error) {
	bn := models.NewBayesianNetwork()
	if err := readXMLNativeWith(r, &realBuilder{bn: bn}); err != nil {
		return nil, err
	}
	return bn, nil
}

// readXMLNativeWith is the testable implementation of ReadXMLNative.
// Accepts a bnBuilder interface for mock injection.
func readXMLNativeWith(r io.Reader, builder bnBuilder) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("readwrite: error reading datascience XML: %w", err)
	}

	var net pgmgoNetwork
	if err := xml.Unmarshal(data, &net); err != nil {
		return fmt.Errorf("readwrite: error parsing datascience XML: %w", err)
	}

	type varInfo struct {
		card   int
		states []string
	}
	varMap := make(map[string]*varInfo)

	// Add nodes.
	for _, node := range net.Nodes.Nodes {
		name := node.Name
		if err := builder.AddNode(name); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
		var states []string
		for _, s := range strings.Split(node.States, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				states = append(states, s)
			}
		}
		if len(states) > 0 {
			if err := builder.SetStates(name, states); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
		}
		varMap[name] = &varInfo{card: len(states), states: states}
	}

	// Add edges.
	for _, edge := range net.Edges.Edges {
		if err := builder.AddEdge(edge.From, edge.To); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("readwrite: %w", err)
			}
		}
	}

	// Add CPDs.
	for _, xc := range net.CPDs.CPDs {
		child := xc.Variable
		childCard := xc.Card

		var parents []string
		var evidenceCard []int

		if strings.TrimSpace(xc.Evidence) != "" {
			for _, p := range strings.Fields(xc.Evidence) {
				p = strings.TrimSpace(p)
				if p != "" {
					parents = append(parents, p)
				}
			}
		}
		if strings.TrimSpace(xc.EvidenceCard) != "" {
			for _, ec := range strings.Fields(xc.EvidenceCard) {
				ec = strings.TrimSpace(ec)
				if ec == "" {
					continue
				}
				v, err := strconv.Atoi(ec)
				if err != nil {
					return fmt.Errorf("readwrite: invalid evidence_card %q: %w", ec, err)
				}
				evidenceCard = append(evidenceCard, v)
			}
		}

		if len(parents) != len(evidenceCard) {
			return fmt.Errorf("readwrite: CPD for %q: evidence count %d != evidence_card count %d",
				child, len(parents), len(evidenceCard))
		}

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		// Parse values (space-separated).
		vals, err := xmlbifParseFloats(xc.Values)
		if err != nil {
			return fmt.Errorf("readwrite: error parsing values for %q: %w", child, err)
		}

		expectedLen := childCard * numParentConfigs
		if len(vals) != expectedLen {
			return fmt.Errorf("readwrite: CPD for %q has %d values, expected %d",
				child, len(vals), expectedLen)
		}

		// Values are stored flat: parent configs outer, child states inner.
		values := make([][]float64, childCard)
		for cs := 0; cs < childCard; cs++ {
			values[cs] = make([]float64, numParentConfigs)
		}

		idx := 0
		for pc := 0; pc < numParentConfigs; pc++ {
			for cs := 0; cs < childCard; cs++ {
				values[cs][pc] = vals[idx]
				idx++
			}
		}

		cpd, err := factors.NewTabularCPD(child, childCard, values, parents, evidenceCard)
		if err != nil {
			return fmt.Errorf("readwrite: failed to create CPD for %q: %w", child, err)
		}
		if err := builder.AddCPD(cpd); err != nil {
			return fmt.Errorf("readwrite: %w", err)
		}
	}

	return nil
}

// WriteXMLNative serializes a BayesianNetwork to datascience-native XML format.
func WriteXMLNative(w io.Writer, bn *models.BayesianNetwork) error {
	nodes := bn.Nodes()

	var xmlNodes []pgmgoNode
	for _, node := range nodes {
		states := bn.GetStates(node)
		if len(states) == 0 {
			return fmt.Errorf("readwrite: variable %q has no state names", node)
		}
		xmlNodes = append(xmlNodes, pgmgoNode{
			Name:   node,
			States: strings.Join(states, ","),
		})
	}

	edges := bn.Edges()
	xmlEdges := make([]pgmgoEdge, len(edges))
	for i, e := range edges {
		xmlEdges[i] = pgmgoEdge{From: e[0], To: e[1]}
	}

	var xmlCPDs []pgmgoCPD
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

		// Build values string: parent configs outer, child states inner.
		var parts []string
		for pc := 0; pc < numParentConfigs; pc++ {
			for cs := 0; cs < childCard; cs++ {
				parts = append(parts, formatFloat(data[cs*numParentConfigs+pc]))
			}
		}

		xc := pgmgoCPD{
			Variable: node,
			Card:     childCard,
			Values:   strings.Join(parts, " "),
		}
		if len(evidence) > 0 {
			xc.Evidence = strings.Join(evidence, " ")
			ecParts := make([]string, len(evidenceCard))
			for i, ec := range evidenceCard {
				ecParts[i] = strconv.Itoa(ec)
			}
			xc.EvidenceCard = strings.Join(ecParts, " ")
		}
		xmlCPDs = append(xmlCPDs, xc)
	}

	net := pgmgoNetwork{
		Name:  "network",
		Nodes: pgmgoNodes{Nodes: xmlNodes},
		Edges: pgmgoEdges{Edges: xmlEdges},
		CPDs:  pgmgoCPDs{CPDs: xmlCPDs},
	}

	if _, err := fmt.Fprint(w, xml.Header); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(net); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}
	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	return nil
}
