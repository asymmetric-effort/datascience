package learning

import (
	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ConstraintBasedEstimator is a structure learning estimator that uses
// conditional independence tests to learn the network structure. This is a
// convenience wrapper around the PC algorithm that matches pgmpy's deprecated
// ConstraintBasedEstimator class.
//
// Deprecated: Use PCAlgorithm directly for new code. This type is retained for
// API compatibility with pgmpy.
type ConstraintBasedEstimator struct {
	pc *PCAlgorithm
}

// NewConstraintBasedEstimator creates a new ConstraintBasedEstimator.
// Parameters match those of NewPC.
func NewConstraintBasedEstimator(data *tabgo.DataFrame, ciTest CITestFunc, significance float64, opts ...PCOption) *ConstraintBasedEstimator {
	return &ConstraintBasedEstimator{
		pc: NewPC(data, ciTest, significance, opts...),
	}
}

// Estimate runs the constraint-based structure learning (PC algorithm) and
// returns the learned PDAG.
func (cbe *ConstraintBasedEstimator) Estimate() (*graphgo.PDAG, error) {
	return cbe.pc.Estimate()
}

// EstimateBN runs the constraint-based structure learning and returns a
// BayesianNetwork (DAG) by orienting remaining undirected edges.
func (cbe *ConstraintBasedEstimator) EstimateBN() (*models.BayesianNetwork, error) {
	return cbe.pc.EstimateBN()
}

// BuildSkeleton runs only the skeleton discovery phase and returns the
// undirected skeleton along with separating sets.
func (cbe *ConstraintBasedEstimator) BuildSkeleton() (*graphgo.PDAG, map[[2]string][]string, error) {
	return cbe.pc.BuildSkeleton()
}
