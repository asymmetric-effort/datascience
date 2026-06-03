package readwrite

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ReadCSVStructure reads an edge-list CSV with "from,to" columns and creates a
// BayesianNetwork containing only structure (nodes and edges, no CPDs).
// The first row must be a header with fields "from" and "to" (case-insensitive).
func ReadCSVStructure(r io.Reader) (*models.BayesianNetwork, error) {
	bn := models.NewBayesianNetwork()
	if err := readCSVStructureWith(r, &realBuilder{bn: bn}); err != nil {
		return nil, err
	}
	return bn, nil
}

// readCSVStructureWith is the testable implementation of ReadCSVStructure.
// Accepts a bnBuilder interface for mock injection.
func readCSVStructureWith(r io.Reader, builder bnBuilder) error {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("readwrite: error reading CSV: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("readwrite: empty CSV file")
	}

	// Parse header.
	header := records[0]
	if len(header) < 2 {
		return fmt.Errorf("readwrite: CSV header must have at least 2 columns, got %d", len(header))
	}

	fromCol, toCol := -1, -1
	for i, h := range header {
		switch strings.TrimSpace(strings.ToLower(h)) {
		case "from":
			fromCol = i
		case "to":
			toCol = i
		}
	}
	if fromCol < 0 || toCol < 0 {
		return fmt.Errorf("readwrite: CSV header must contain 'from' and 'to' columns, got %v", header)
	}

	added := make(map[string]bool)

	for _, row := range records[1:] {
		if len(row) <= fromCol || len(row) <= toCol {
			continue
		}
		from := strings.TrimSpace(row[fromCol])
		to := strings.TrimSpace(row[toCol])
		if from == "" || to == "" {
			continue
		}

		if !added[from] {
			if err := builder.AddNode(from); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
			added[from] = true
		}
		if !added[to] {
			if err := builder.AddNode(to); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
			added[to] = true
		}
		if err := builder.AddEdge(from, to); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("readwrite: %w", err)
			}
		}
	}

	return nil
}

// WriteCSVStructure writes the edge list of a BayesianNetwork as CSV with
// columns "from,to".
func WriteCSVStructure(w io.Writer, bn *models.BayesianNetwork) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write([]string{"from", "to"}); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	for _, edge := range bn.Edges() {
		if err := writer.Write([]string{edge[0], edge[1]}); err != nil {
			return fmt.Errorf("readwrite: write error: %w", err)
		}
	}

	return nil
}

// ReadCSVCPD reads a CPD from a CSV table.
//
// Format:
//   - Row 0 (header): first cell is the variable name, remaining cells are
//     parent configuration labels (or empty for unconditional).
//   - Row 1..N: first cell is the state name, remaining cells are probabilities.
//
// For an unconditional CPD (no parents), the header has one cell and each data
// row has two cells (state, probability).
//
// For a conditional CPD, parent configurations are encoded as
// "Parent1=state,Parent2=state" in the header cells.
func ReadCSVCPD(r io.Reader) (*factors.TabularCPD, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("readwrite: error reading CSV CPD: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("readwrite: CSV CPD must have at least a header and one data row")
	}

	header := records[0]
	variable := strings.TrimSpace(header[0])
	if variable == "" {
		return nil, fmt.Errorf("readwrite: CSV CPD header first cell must be the variable name")
	}

	numParentConfigs := len(header) - 1
	if numParentConfigs < 1 {
		return nil, fmt.Errorf("readwrite: CSV CPD must have at least one value column")
	}

	// Parse parent info from header columns (if any).
	var parents []string
	var evidenceCard []int

	// Check if this is conditional (header cells contain "Parent=state" patterns).
	if numParentConfigs > 1 || (numParentConfigs == 1 && strings.Contains(header[1], "=")) {
		// Parse parent configurations from header.
		parentConfigs, pNames, err := csvParseCPDHeader(header[1:])
		if err != nil {
			return nil, err
		}
		parents = pNames

		// Determine cardinalities from the configs.
		parentStatesSeen := make(map[string]map[string]bool)
		for _, p := range parents {
			parentStatesSeen[p] = make(map[string]bool)
		}
		for _, cfg := range parentConfigs {
			for p, s := range cfg {
				parentStatesSeen[p][s] = true
			}
		}
		for _, p := range parents {
			evidenceCard = append(evidenceCard, len(parentStatesSeen[p]))
		}
		// Verify numParentConfigs matches expected.
		expected := 1
		for _, ec := range evidenceCard {
			expected *= ec
		}
		if numParentConfigs != expected {
			return nil, fmt.Errorf("readwrite: CSV CPD has %d parent configs but expected %d from cardinalities",
				numParentConfigs, expected)
		}
	}

	// Parse data rows.
	childCard := len(records) - 1
	values := make([][]float64, childCard)
	for i := 0; i < childCard; i++ {
		row := records[i+1]
		if len(row) < numParentConfigs+1 {
			return nil, fmt.Errorf("readwrite: CSV CPD data row %d has %d cells, expected %d",
				i, len(row), numParentConfigs+1)
		}
		values[i] = make([]float64, numParentConfigs)
		for j := 0; j < numParentConfigs; j++ {
			v, err := strconv.ParseFloat(strings.TrimSpace(row[j+1]), 64)
			if err != nil {
				return nil, fmt.Errorf("readwrite: CSV CPD invalid value at row %d col %d: %w", i+1, j+1, err)
			}
			values[i][j] = v
		}
	}

	cpd, err := factors.NewTabularCPD(variable, childCard, values, parents, evidenceCard)
	if err != nil {
		return nil, fmt.Errorf("readwrite: failed to create CPD from CSV: %w", err)
	}
	return cpd, nil
}

// WriteCSVCPD writes a TabularCPD as a CSV table.
//
// The output has a header row with the variable name in the first cell, followed
// by parent configuration labels (or a single empty-label column for unconditional).
// Each subsequent row is a child state followed by its probability values.
func WriteCSVCPD(w io.Writer, cpd *factors.TabularCPD) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	evidence := cpd.Evidence()
	evidenceCard := cpd.EvidenceCard()
	childCard := cpd.VariableCard()

	numParentConfigs := 1
	for _, ec := range evidenceCard {
		numParentConfigs *= ec
	}

	data := cpd.ToFactor().Values().Data()

	// Build header.
	header := []string{cpd.Variable()}
	if len(evidence) == 0 {
		header = append(header, "P")
	} else {
		for pc := 0; pc < numParentConfigs; pc++ {
			label := csvParentConfigLabel(pc, evidence, evidenceCard)
			header = append(header, label)
		}
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	// Write data rows.
	for cs := 0; cs < childCard; cs++ {
		row := []string{fmt.Sprintf("s%d", cs)}
		for pc := 0; pc < numParentConfigs; pc++ {
			row = append(row, formatFloat(data[cs*numParentConfigs+pc]))
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("readwrite: write error: %w", err)
		}
	}

	return nil
}

// csvParentConfigLabel builds a label like "A=s0,B=s1" for a given parent config index.
func csvParentConfigLabel(pc int, evidence []string, evidenceCard []int) string {
	indices := make([]int, len(evidence))
	rem := pc
	for i := len(evidence) - 1; i >= 0; i-- {
		indices[i] = rem % evidenceCard[i]
		rem /= evidenceCard[i]
	}

	var parts []string
	for i, ev := range evidence {
		parts = append(parts, fmt.Sprintf("%s=s%d", ev, indices[i]))
	}
	return strings.Join(parts, ",")
}

// csvParseCPDHeader parses parent configuration labels from header cells.
// Returns the list of configs (each a map[parent]state) and the ordered parent names.
func csvParseCPDHeader(cells []string) ([]map[string]string, []string, error) {
	var configs []map[string]string
	var parentOrder []string
	parentSeen := make(map[string]bool)

	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		cfg := make(map[string]string)
		for _, part := range strings.Split(cell, ",") {
			part = strings.TrimSpace(part)
			eqIdx := strings.Index(part, "=")
			if eqIdx < 0 {
				return nil, nil, fmt.Errorf("readwrite: CSV CPD header cell %q missing '=' in parent config", cell)
			}
			parent := strings.TrimSpace(part[:eqIdx])
			state := strings.TrimSpace(part[eqIdx+1:])
			cfg[parent] = state
			if !parentSeen[parent] {
				parentSeen[parent] = true
				parentOrder = append(parentOrder, parent)
			}
		}
		configs = append(configs, cfg)
	}

	return configs, parentOrder, nil
}
