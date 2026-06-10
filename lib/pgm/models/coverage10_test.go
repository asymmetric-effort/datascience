//go:build unit

package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ---------------------------------------------------------------------------
// nbFitImpl: CPD creation failure for class CPD.
// ---------------------------------------------------------------------------
func TestNbFitImpl_ClassCPDFailure(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"C":  tabgo.NewSeries("C", []any{0, 1, 0, 1}),
		"F1": tabgo.NewSeries("F1", []any{0, 1, 0, 1}),
	})
	err := nbFitImpl(nb, df, failingCPDCreator)
	if err == nil || !strings.Contains(err.Error(), "class CPD") {
		t.Fatalf("expected class CPD failure, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// nbFitImpl: CPD creation failure for feature CPD (class succeeds, feature fails).
// ---------------------------------------------------------------------------
func TestNbFitImpl_FeatureCPDFailure(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"C":  tabgo.NewSeries("C", []any{0, 1, 0, 1}),
		"F1": tabgo.NewSeries("F1", []any{0, 1, 0, 1}),
	})
	// Creator that succeeds for class but fails for features.
	callCount := 0
	featureFailCreator := func(variable string, variableCard int, values [][]float64, evidence []string, evidenceCard []int) (*tabularCPDType, error) {
		callCount++
		if callCount > 1 { // First call is class CPD (success), second is feature (fail).
			return nil, failingCPDCreatorErr(variable)
		}
		return defaultCPDCreator(variable, variableCard, values, evidence, evidenceCard)
	}
	err := nbFitImpl(nb, df, featureFailCreator)
	if err == nil || !strings.Contains(err.Error(), "CPD") {
		t.Fatalf("expected CPD failure, got: %v", err)
	}
}

// We need a type alias and helper for the failing creator.
type tabularCPDType = factors.TabularCPD

func failingCPDCreatorErr(variable string) error {
	return fmt.Errorf("injected CPD creation failure for %q", variable)
}

// ---------------------------------------------------------------------------
// nbFitImpl: nil data and empty data.
// ---------------------------------------------------------------------------
func TestNbFitImpl_NilData(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	err := nbFitImpl(nb, nil, defaultCPDCreator)
	if err == nil || !strings.Contains(err.Error(), "nil") {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestNbFitImpl_EmptyData(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{})
	err := nbFitImpl(nb, df, defaultCPDCreator)
	if err == nil {
		t.Fatal("expected error for empty data")
	}
}

// ---------------------------------------------------------------------------
// nbFitImpl: negative class value.
// ---------------------------------------------------------------------------
func TestNbFitImpl_NegativeClassValue(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"C":  tabgo.NewSeries("C", []any{-1}),
		"F1": tabgo.NewSeries("F1", []any{0}),
	})
	err := nbFitImpl(nb, df, defaultCPDCreator)
	if err == nil || !strings.Contains(err.Error(), "negative") {
		t.Fatalf("expected negative error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// nbFitImpl: negative feature value.
// ---------------------------------------------------------------------------
func TestNbFitImpl_NegativeFeatureValue(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"C":  tabgo.NewSeries("C", []any{0, 1}),
		"F1": tabgo.NewSeries("F1", []any{0, -1}),
	})
	err := nbFitImpl(nb, df, defaultCPDCreator)
	if err == nil || !strings.Contains(err.Error(), "negative") {
		t.Fatalf("expected negative error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// nbFitImpl: zero class count (uniform distribution fallback).
// ---------------------------------------------------------------------------
func TestNbFitImpl_ZeroClassCount(t *testing.T) {
	nb, _ := NewNaiveBayes("C", []string{"F1"})
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"C":  tabgo.NewSeries("C", []any{0, 0, 0, 2}),
		"F1": tabgo.NewSeries("F1", []any{0, 1, 0, 0}),
	})
	err := nbFitImpl(nb, df, defaultCPDCreator)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
