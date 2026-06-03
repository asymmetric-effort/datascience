package readwrite

import (
	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// bnBuilder abstracts BayesianNetwork construction to enable
// testing of defensive error paths in format readers.
type bnBuilder interface {
	AddNode(name string) error
	AddEdge(from, to string) error
	SetStates(node string, states []string) error
	AddCPD(cpd *factors.TabularCPD) error
}

// realBuilder wraps a real BayesianNetwork as a bnBuilder.
type realBuilder struct {
	bn *models.BayesianNetwork
}

func (r *realBuilder) AddNode(name string) error            { return r.bn.AddNode(name) }
func (r *realBuilder) AddEdge(from, to string) error        { return r.bn.AddEdge(from, to) }
func (r *realBuilder) SetStates(n string, s []string) error { return r.bn.SetStates(n, s) }
func (r *realBuilder) AddCPD(cpd *factors.TabularCPD) error { return r.bn.AddCPD(cpd) }
