//go:build unit

package tabgo

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadWriteCSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	// Write a CSV manually
	content := "name,age,city\nAlice,30,NYC\nBob,25,LA\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Read it
	df, err := ReadCSV(path)
	if err != nil {
		t.Fatalf("ReadCSV: %v", err)
	}
	if df.Len() != 2 {
		t.Fatalf("Len = %d, want 2", df.Len())
	}
	cols := df.Columns()
	if !reflect.DeepEqual(cols, []string{"name", "age", "city"}) {
		t.Fatalf("Columns = %v", cols)
	}
	names := df.Column("name").Values()
	if !reflect.DeepEqual(names, []any{"Alice", "Bob"}) {
		t.Fatalf("name values = %v", names)
	}
	ages := df.Column("age").Values()
	if !reflect.DeepEqual(ages, []any{"30", "25"}) {
		t.Fatalf("age values = %v (CSV stores strings)", ages)
	}
}

func TestWriteCSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")

	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{"hello", 42},
			{"world", 99},
		},
	)
	if err := WriteCSV(df, path); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}

	// Read back and verify
	df2, err := ReadCSV(path)
	if err != nil {
		t.Fatalf("ReadCSV after write: %v", err)
	}
	if df2.Len() != 2 {
		t.Fatalf("Len = %d", df2.Len())
	}
	aVals := df2.Column("a").Values()
	if !reflect.DeepEqual(aVals, []any{"hello", "world"}) {
		t.Fatalf("a = %v", aVals)
	}
	bVals := df2.Column("b").Values()
	if !reflect.DeepEqual(bVals, []any{"42", "99"}) {
		t.Fatalf("b = %v", bVals)
	}
}

func TestReadCSVEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	df, err := ReadCSV(path)
	if err != nil {
		t.Fatalf("ReadCSV empty: %v", err)
	}
	if df.Len() != 0 {
		t.Fatalf("empty CSV Len = %d", df.Len())
	}
}

func TestReadCSVHeaderOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "header.csv")
	if err := os.WriteFile(path, []byte("a,b,c\n"), 0644); err != nil {
		t.Fatal(err)
	}
	df, err := ReadCSV(path)
	if err != nil {
		t.Fatalf("ReadCSV header-only: %v", err)
	}
	if df.Len() != 0 {
		t.Fatalf("header-only Len = %d", df.Len())
	}
	if !reflect.DeepEqual(df.Columns(), []string{"a", "b", "c"}) {
		t.Fatalf("Columns = %v", df.Columns())
	}
}

func TestReadCSVNotFound(t *testing.T) {
	_, err := ReadCSV("/nonexistent/path.csv")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWriteCSVNilValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nil.csv")
	df := NewDataFrameFromRows([]string{"x"}, [][]any{{nil}})
	if err := WriteCSV(df, path); err != nil {
		t.Fatalf("WriteCSV with nil: %v", err)
	}
	data, _ := os.ReadFile(path)
	// should contain empty string for nil
	if string(data) != "x\n\n" {
		t.Fatalf("nil output = %q", string(data))
	}
}

func TestRoundTripFloat(t *testing.T) {
	// CSV stores everything as strings; verify Float64 conversion works on read data
	dir := t.TempDir()
	path := filepath.Join(dir, "float.csv")
	if err := os.WriteFile(path, []byte("val\n1.5\n2.7\n3.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	df, err := ReadCSV(path)
	if err != nil {
		t.Fatal(err)
	}
	floats := df.Column("val").Float64()
	want := []float64{1.5, 2.7, 3.0}
	for i := range floats {
		diff := floats[i] - want[i]
		if diff < -0.001 || diff > 0.001 {
			t.Fatalf("[%d] = %f, want %f", i, floats[i], want[i])
		}
	}
}
