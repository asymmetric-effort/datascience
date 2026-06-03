package learning

import (
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// BaseEstimator is the common base type for all parameter estimation algorithms.
// It holds references to the BayesianNetwork being estimated and the observed
// data. Concrete estimators (MLE, Bayesian, EM, etc.) embed this type.
//
// This mirrors pgmpy's BaseEstimator base class.
type BaseEstimator struct {
	Model *models.BayesianNetwork
	Data  *tabgo.DataFrame
}

// NewBaseEstimator creates a new BaseEstimator with the given model and data.
func NewBaseEstimator(model *models.BayesianNetwork, data *tabgo.DataFrame) BaseEstimator {
	return BaseEstimator{
		Model: model,
		Data:  data,
	}
}

// GetModel returns the BayesianNetwork associated with this estimator.
func (b *BaseEstimator) GetModel() *models.BayesianNetwork {
	return b.Model
}

// GetData returns the DataFrame associated with this estimator.
func (b *BaseEstimator) GetData() *tabgo.DataFrame {
	return b.Data
}

// StructureEstimator is the common base type for all structure learning
// algorithms. It holds a reference to the data and provides shared utility
// methods. Concrete structure learners (PC, HillClimb, GES, etc.) can embed
// this type.
//
// This mirrors pgmpy's StructureEstimator base class.
type StructureEstimator struct {
	Data *tabgo.DataFrame
}

// NewStructureEstimator creates a new StructureEstimator with the given data.
func NewStructureEstimator(data *tabgo.DataFrame) StructureEstimator {
	return StructureEstimator{
		Data: data,
	}
}

// GetData returns the DataFrame associated with this structure estimator.
func (s *StructureEstimator) GetData() *tabgo.DataFrame {
	return s.Data
}

// Variables returns the column names (variable names) from the data.
func (s *StructureEstimator) Variables() []string {
	if s.Data == nil {
		return nil
	}
	cols := s.Data.Columns()
	result := make([]string, len(cols))
	copy(result, cols)
	return result
}
