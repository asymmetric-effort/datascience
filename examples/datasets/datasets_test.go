//go:build unit

package datasets

import (
	"testing"
)

// datasetSpec defines expected row and column counts for each dataset.
type datasetSpec struct {
	name string
	rows int
	cols int
}

var allDatasets = []datasetSpec{
	// Original 8
	{"asia", 1000, 8},
	{"alarm", 1000, 5},
	{"sachs", 500, 11},
	{"cancer", 1000, 5},
	{"student", 1000, 5},
	{"sprinkler", 1000, 4},
	{"survey", 500, 5},
	{"titanic", 800, 4},
	// Well-known ML
	{"adult", 500, 11},
	{"pima_diabetes", 300, 9},
	{"iris", 150, 5},
	{"wine", 300, 12},
	{"heart", 300, 14},
	{"boston", 300, 14},
	{"breast_cancer", 300, 10},
	// BN-specific
	{"earthquake", 500, 5},
	{"child", 500, 20},
	{"insurance", 500, 27},
	{"water", 500, 32},
	{"mildew", 500, 35},
	{"hailfinder", 500, 56},
	{"hepar2", 500, 40},
	{"lucas", 500, 12},
	{"andes", 500, 30},
	{"munin", 500, 30},
	{"barley", 500, 30},
	{"win95pts", 500, 30},
	// UCI
	{"ecoli", 336, 8},
	{"glass", 214, 10},
	{"zoo", 101, 17},
	{"mushroom", 500, 23},
	{"nursery", 500, 9},
	{"car_evaluation", 400, 7},
	{"balance_scale", 300, 5},
	{"monks", 300, 7},
	{"tic_tac_toe", 400, 10},
	{"vote", 435, 17},
	{"credit_approval", 300, 16},
	{"hepatitis", 155, 20},
	{"automobile", 205, 26},
}

func TestAllDatasets(t *testing.T) {
	for _, ds := range allDatasets {
		t.Run(ds.name, func(t *testing.T) {
			df, err := Load(ds.name)
			if err != nil {
				t.Fatalf("Load(%q): %v", ds.name, err)
			}
			if df.Len() != ds.rows {
				t.Errorf("%s: expected %d rows, got %d", ds.name, ds.rows, df.Len())
			}
			cols := df.Columns()
			if len(cols) != ds.cols {
				t.Errorf("%s: expected %d columns, got %d", ds.name, ds.cols, len(cols))
			}
		})
	}
}

func TestList(t *testing.T) {
	names := List()
	if len(names) != len(allDatasets) {
		t.Fatalf("List: expected %d datasets, got %d", len(allDatasets), len(names))
	}
	// Verify sorted order
	for i := 1; i < len(names); i++ {
		if names[i] <= names[i-1] {
			t.Errorf("List not sorted: %q <= %q at index %d", names[i], names[i-1], i)
		}
	}
}

func TestLoad(t *testing.T) {
	for _, name := range List() {
		df, err := Load(name)
		if err != nil {
			t.Errorf("Load(%q): %v", name, err)
			continue
		}
		if df.Len() == 0 {
			t.Errorf("Load(%q): got 0 rows", name)
		}
	}
}

func TestLoadUnknown(t *testing.T) {
	_, err := Load("nonexistent")
	if err == nil {
		t.Error("Load(nonexistent): expected error, got nil")
	}
}

// Preserve original individual test functions for backward compatibility.

func TestAsia(t *testing.T) {
	df, err := Asia()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 1000 {
		t.Errorf("Asia: expected 1000 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 8 {
		t.Errorf("Asia: expected 8 columns, got %d", len(cols))
	}
}

func TestAlarm(t *testing.T) {
	df, err := Alarm()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 1000 {
		t.Errorf("Alarm: expected 1000 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 5 {
		t.Errorf("Alarm: expected 5 columns, got %d", len(cols))
	}
}

func TestSachs(t *testing.T) {
	df, err := Sachs()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 500 {
		t.Errorf("Sachs: expected 500 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 11 {
		t.Errorf("Sachs: expected 11 columns, got %d", len(cols))
	}
}

func TestCancer(t *testing.T) {
	df, err := Cancer()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 1000 {
		t.Errorf("Cancer: expected 1000 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 5 {
		t.Errorf("Cancer: expected 5 columns, got %d", len(cols))
	}
}

func TestStudent(t *testing.T) {
	df, err := Student()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 1000 {
		t.Errorf("Student: expected 1000 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 5 {
		t.Errorf("Student: expected 5 columns, got %d", len(cols))
	}
}

func TestSprinkler(t *testing.T) {
	df, err := Sprinkler()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 1000 {
		t.Errorf("Sprinkler: expected 1000 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 4 {
		t.Errorf("Sprinkler: expected 4 columns, got %d", len(cols))
	}
}

func TestSurvey(t *testing.T) {
	df, err := Survey()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 500 {
		t.Errorf("Survey: expected 500 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 5 {
		t.Errorf("Survey: expected 5 columns, got %d", len(cols))
	}
}

func TestTitanic(t *testing.T) {
	df, err := Titanic()
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 800 {
		t.Errorf("Titanic: expected 800 rows, got %d", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 4 {
		t.Errorf("Titanic: expected 4 columns, got %d", len(cols))
	}
}
