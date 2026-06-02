//go:build unit

package tabgo

import (
	"strings"
	"testing"
)

func TestToJSON(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name", "age"},
		[][]any{
			{"Alice", 30},
			{"Bob", 25},
		},
	)
	jsonStr, err := ToJSON(df)
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	if !strings.Contains(jsonStr, "Alice") {
		t.Error("ToJSON missing Alice")
	}
	if !strings.Contains(jsonStr, "Bob") {
		t.Error("ToJSON missing Bob")
	}
}

func TestReadJSON(t *testing.T) {
	data := `[{"name":"Alice","age":30},{"name":"Bob","age":25}]`
	df, err := ReadJSON(data)
	if err != nil {
		t.Fatalf("ReadJSON error: %v", err)
	}
	if df.Len() != 2 {
		t.Errorf("ReadJSON rows = %d, want 2", df.Len())
	}
	cols := df.Columns()
	if len(cols) != 2 {
		t.Errorf("ReadJSON cols = %d, want 2", len(cols))
	}
}

func TestReadJSONEmpty(t *testing.T) {
	df, err := ReadJSON("[]")
	if err != nil {
		t.Fatalf("ReadJSON empty error: %v", err)
	}
	if df.Len() != 0 {
		t.Errorf("ReadJSON empty rows = %d, want 0", df.Len())
	}
}

func TestReadJSONInvalid(t *testing.T) {
	_, err := ReadJSON("not json")
	if err == nil {
		t.Error("ReadJSON should fail on invalid json")
	}
}

func TestToJSONRoundTrip(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x", "y"},
		[][]any{
			{1.0, 2.0},
			{3.0, 4.0},
		},
	)
	jsonStr, err := ToJSON(df)
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	df2, err := ReadJSON(jsonStr)
	if err != nil {
		t.Fatalf("ReadJSON error: %v", err)
	}
	if df2.Len() != 2 {
		t.Errorf("Round trip rows = %d, want 2", df2.Len())
	}
}

func TestToHTML(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"col"},
		[][]any{{"val"}},
	)
	html := ToHTML(df)
	if !strings.Contains(html, "<table>") {
		t.Error("ToHTML missing <table>")
	}
	if !strings.Contains(html, "<th>col</th>") {
		t.Error("ToHTML missing header")
	}
	if !strings.Contains(html, "<td>val</td>") {
		t.Error("ToHTML missing data")
	}
}

func TestToHTMLEscaping(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{"<b>bold</b>"}},
	)
	html := ToHTML(df)
	if strings.Contains(html, "<b>bold</b>") {
		t.Error("ToHTML should escape HTML in values")
	}
	if !strings.Contains(html, "&lt;b&gt;") {
		t.Error("ToHTML should have escaped angle brackets")
	}
}

func TestToDict(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, 2},
			{3, 4},
		},
	)
	d := ToDict(df)
	if len(d) != 2 {
		t.Errorf("ToDict cols = %d, want 2", len(d))
	}
	if len(d["a"]) != 2 {
		t.Errorf("ToDict['a'] len = %d, want 2", len(d["a"]))
	}
	if d["a"][0] != 1 || d["a"][1] != 3 {
		t.Errorf("ToDict['a'] = %v, want [1, 3]", d["a"])
	}
}

func TestToRecords(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name", "val"},
		[][]any{
			{"x", 1},
			{"y", 2},
		},
	)
	records := ToRecords(df)
	if len(records) != 2 {
		t.Fatalf("ToRecords len = %d, want 2", len(records))
	}
	if records[0]["name"] != "x" {
		t.Errorf("record[0]['name'] = %v, want x", records[0]["name"])
	}
	if records[1]["val"] != 2 {
		t.Errorf("record[1]['val'] = %v, want 2", records[1]["val"])
	}
}
