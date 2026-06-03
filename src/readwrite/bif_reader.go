package readwrite

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// bifVarInfo holds parsed variable metadata.
type bifVarInfo struct {
	name   string
	card   int
	states []string
}

// ReadBIF parses a BIF (Bayesian Interchange Format) file and returns a
// BayesianNetwork. It handles network, variable, and probability blocks.
// Lines containing // comments have the comment portion stripped.
func ReadBIF(r io.Reader) (*models.BayesianNetwork, error) {
	bn := models.NewBayesianNetwork()
	if err := readBIFWith(r, &realBuilder{bn: bn}); err != nil {
		return nil, err
	}
	return bn, nil
}

// readBIFWith is the testable implementation of ReadBIF. Accepts a bnBuilder
// interface for mock injection.
func readBIFWith(r io.Reader, builder bnBuilder) error {
	lines, err := readBIFLines(r)
	if err != nil {
		return err
	}

	varMap := make(map[string]*bifVarInfo)

	i := 0
	for i < len(lines) {
		tokens := strings.Fields(lines[i])
		if len(tokens) == 0 {
			i++
			continue
		}

		switch tokens[0] {
		case "network":
			i = bifSkipBlock(lines, i)

		case "variable":
			if len(tokens) < 2 {
				return fmt.Errorf("readwrite: malformed variable declaration")
			}
			name := strings.TrimRight(tokens[1], "{")
			if name == "" {
				return fmt.Errorf("readwrite: malformed variable declaration: missing name")
			}
			blockContent, end := bifCollectBlock(lines, i)
			i = end

			states, err := bifParseVariableBlock(name, blockContent)
			if err != nil {
				return err
			}

			if err := builder.AddNode(name); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
			if err := builder.SetStates(name, states); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}
			varMap[name] = &bifVarInfo{name: name, card: len(states), states: states}

		case "probability":
			// Find the probability header line (the line containing "probability").
			headerLine := lines[i]
			blockContent, end := bifCollectBlock(lines, i)
			i = end

			child, parents, err := bifParseProbHeader(headerLine)
			if err != nil {
				return err
			}

			childInfo := varMap[child]
			if childInfo == nil {
				return fmt.Errorf("readwrite: probability references unknown variable %q", child)
			}

			// Add edges.
			for _, p := range parents {
				if err := builder.AddEdge(p, child); err != nil {
					if !strings.Contains(err.Error(), "already exists") {
						return fmt.Errorf("readwrite: %w", err)
					}
				}
			}

			// Gather parent info.
			var parentInfos []*bifVarInfo
			for _, p := range parents {
				pi := varMap[p]
				if pi == nil {
					return fmt.Errorf("readwrite: probability references unknown parent %q", p)
				}
				parentInfos = append(parentInfos, pi)
			}

			cpd, err := bifParseProbBlock(childInfo, parents, parentInfos, blockContent)
			if err != nil {
				return err
			}
			if err := builder.AddCPD(cpd); err != nil {
				return fmt.Errorf("readwrite: %w", err)
			}

		default:
			i++
		}
	}

	return nil
}

// readBIFLines reads all lines, strips comments, and returns non-empty trimmed lines.
func readBIFLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		line := bifStripComment(scanner.Text())
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("readwrite: error reading BIF: %w", err)
	}
	return lines, nil
}

// bifStripComment removes // comments from a line.
func bifStripComment(line string) string {
	if idx := strings.Index(line, "//"); idx >= 0 {
		return line[:idx]
	}
	return line
}

// bifSkipBlock advances past a { ... } block starting at line i.
func bifSkipBlock(lines []string, i int) int {
	depth := 0
	for i < len(lines) {
		depth += strings.Count(lines[i], "{") - strings.Count(lines[i], "}")
		i++
		if depth <= 0 {
			break
		}
	}
	return i
}

// bifCollectBlock returns the content lines inside the { } block starting at
// line i, and the index after the closing }.
func bifCollectBlock(lines []string, start int) ([]string, int) {
	// Join all lines from start until we close the brace, then extract inner content.
	depth := 0
	var raw []string
	i := start
	for i < len(lines) {
		depth += strings.Count(lines[i], "{") - strings.Count(lines[i], "}")
		raw = append(raw, lines[i])
		i++
		if depth <= 0 {
			break
		}
	}

	// Join and extract content between first { and last }.
	joined := strings.Join(raw, "\n")
	openIdx := strings.Index(joined, "{")
	closeIdx := strings.LastIndex(joined, "}")
	if openIdx < 0 || closeIdx <= openIdx {
		return nil, i
	}
	inner := joined[openIdx+1 : closeIdx]

	// Split inner content back to lines.
	var content []string
	for _, line := range strings.Split(inner, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			content = append(content, line)
		}
	}
	return content, i
}

// bifParseVariableBlock parses lines inside a variable { } block.
func bifParseVariableBlock(name string, blockLines []string) ([]string, error) {
	for _, line := range blockLines {
		if strings.HasPrefix(strings.TrimSpace(line), "type") {
			return bifParseTypeDecl(name, line)
		}
	}
	return nil, fmt.Errorf("readwrite: no type declaration for variable %q", name)
}

// bifParseTypeDecl parses: type discrete [ N ] { State1, State2, ... };
func bifParseTypeDecl(varName, line string) ([]string, error) {
	line = strings.TrimRight(strings.TrimSpace(line), ";")

	openBrace := strings.Index(line, "{")
	closeBrace := strings.LastIndex(line, "}")
	if openBrace < 0 || closeBrace <= openBrace {
		return nil, fmt.Errorf("readwrite: malformed type declaration for %q", varName)
	}

	stateStr := line[openBrace+1 : closeBrace]
	var states []string
	for _, p := range strings.Split(stateStr, ",") {
		s := strings.TrimSpace(p)
		if s != "" {
			states = append(states, s)
		}
	}
	if len(states) == 0 {
		return nil, fmt.Errorf("readwrite: no states found for variable %q", varName)
	}
	return states, nil
}

// bifParseProbHeader parses: probability ( Child | Parent1, Parent2 ) {
func bifParseProbHeader(line string) (string, []string, error) {
	openParen := strings.Index(line, "(")
	closeParen := strings.LastIndex(line, ")")
	if openParen < 0 || closeParen <= openParen {
		return "", nil, fmt.Errorf("readwrite: malformed probability header: %s", line)
	}

	inner := strings.TrimSpace(line[openParen+1 : closeParen])
	parts := strings.SplitN(inner, "|", 2)
	child := strings.TrimSpace(parts[0])
	if child == "" {
		return "", nil, fmt.Errorf("readwrite: empty variable in probability header")
	}

	var parents []string
	if len(parts) == 2 {
		for _, p := range strings.Split(parts[1], ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				parents = append(parents, p)
			}
		}
	}
	return child, parents, nil
}

// bifParseProbBlock parses the content of a probability { } block into a TabularCPD.
func bifParseProbBlock(child *bifVarInfo, parents []string, parentInfos []*bifVarInfo, blockLines []string) (*factors.TabularCPD, error) {
	numParentConfigs := 1
	var evidenceCard []int
	for _, pi := range parentInfos {
		numParentConfigs *= pi.card
		evidenceCard = append(evidenceCard, pi.card)
	}

	// values[childState][parentConfig] — we'll fill this.
	values := make([][]float64, child.card)
	for i := range values {
		values[i] = make([]float64, numParentConfigs)
	}

	if len(parents) == 0 {
		// Unconditional: "table 0.2, 0.8;"
		for _, line := range blockLines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "table") {
				line = strings.TrimPrefix(line, "table")
				vals, err := bifParseFloats(line)
				if err != nil {
					return nil, fmt.Errorf("readwrite: error parsing table for %q: %w", child.name, err)
				}
				if len(vals) != child.card {
					return nil, fmt.Errorf("readwrite: table for %q has %d values, expected %d",
						child.name, len(vals), child.card)
				}
				for i, v := range vals {
					values[i][0] = v
				}
				break
			}
		}
	} else {
		// Conditional: "(True) 0.01, 0.99;" or "(True, High) 0.1, 0.9;"
		for _, line := range blockLines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "table") {
				// Some BIF files use "table" for conditional too.
				line = strings.TrimPrefix(line, "table")
				vals, err := bifParseFloats(line)
				if err != nil {
					return nil, fmt.Errorf("readwrite: error parsing table for %q: %w", child.name, err)
				}
				if len(vals) != child.card*numParentConfigs {
					return nil, fmt.Errorf("readwrite: table for %q has %d values, expected %d",
						child.name, len(vals), child.card*numParentConfigs)
				}
				// Values are in row-major order: iterate child states, then parent configs.
				idx := 0
				for pc := 0; pc < numParentConfigs; pc++ {
					for cs := 0; cs < child.card; cs++ {
						values[cs][pc] = vals[idx]
						idx++
					}
				}
				break
			}

			if !strings.HasPrefix(line, "(") {
				continue
			}

			// Parse "(State1, State2) val1, val2;"
			closeParen := strings.Index(line, ")")
			if closeParen < 0 {
				return nil, fmt.Errorf("readwrite: malformed conditional line: %s", line)
			}
			stateStr := line[1:closeParen]
			valStr := line[closeParen+1:]

			// Parse parent states.
			var parentStates []string
			for _, s := range strings.Split(stateStr, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					parentStates = append(parentStates, s)
				}
			}
			if len(parentStates) != len(parents) {
				return nil, fmt.Errorf("readwrite: conditional line has %d parent states, expected %d",
					len(parentStates), len(parents))
			}

			// Compute parent config index (row-major over parent variables).
			parentConfig, err := bifParentConfigIndex(parentStates, parentInfos)
			if err != nil {
				return nil, fmt.Errorf("readwrite: %w", err)
			}

			// Parse child values.
			vals, err := bifParseFloats(valStr)
			if err != nil {
				return nil, fmt.Errorf("readwrite: error parsing values: %w", err)
			}
			if len(vals) != child.card {
				return nil, fmt.Errorf("readwrite: conditional line has %d values, expected %d",
					len(vals), child.card)
			}
			for cs := 0; cs < child.card; cs++ {
				values[cs][parentConfig] = vals[cs]
			}
		}
	}

	cpd, err := factors.NewTabularCPD(child.name, child.card, values, parents, evidenceCard)
	if err != nil {
		return nil, fmt.Errorf("readwrite: failed to create CPD for %q: %w", child.name, err)
	}
	return cpd, nil
}

// bifParentConfigIndex computes the flat row-major index for a parent state combination.
func bifParentConfigIndex(parentStates []string, parentInfos []*bifVarInfo) (int, error) {
	idx := 0
	stride := 1
	for i := len(parentInfos) - 1; i >= 0; i-- {
		stateIdx := -1
		for j, s := range parentInfos[i].states {
			if s == parentStates[i] {
				stateIdx = j
				break
			}
		}
		if stateIdx < 0 {
			return 0, fmt.Errorf("unknown state %q for parent %q", parentStates[i], parentInfos[i].name)
		}
		idx += stateIdx * stride
		stride *= parentInfos[i].card
	}
	return idx, nil
}

// bifParseFloats parses comma-separated float values, ignoring trailing semicolons.
func bifParseFloats(s string) ([]float64, error) {
	s = strings.TrimRight(strings.TrimSpace(s), ";")
	var vals []float64
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid probability value %q: %w", p, err)
		}
		vals = append(vals, v)
	}
	return vals, nil
}
