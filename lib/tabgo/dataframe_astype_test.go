//go:build unit

package tabgo

import (
	"testing"
)

func TestDataFrameAstype(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, "10"},
			{2, "20"},
			{3, "30"},
		},
	)
	result := df.Astype("b", "float64")
	vals := result.Column("b").Values()
	if v, ok := vals[0].(float64); !ok || v != 10.0 {
		t.Errorf("Astype: col b[0] = %v (%T), want 10.0", vals[0], vals[0])
	}
}

func TestDataFrameConvertDtypes(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"nums", "bools", "strs"},
		[][]any{
			{"1.5", "true", "hello"},
			{"2.5", "false", "world"},
		},
	)
	result := df.ConvertDtypes()

	// nums should be float64
	numVals := result.Column("nums").Values()
	if v, ok := numVals[0].(float64); !ok || v != 1.5 {
		t.Errorf("ConvertDtypes: nums[0] = %v (%T), want 1.5", numVals[0], numVals[0])
	}

	// bools should be bool
	boolVals := result.Column("bools").Values()
	if v, ok := boolVals[0].(bool); !ok || v != true {
		t.Errorf("ConvertDtypes: bools[0] = %v (%T), want true", boolVals[0], boolVals[0])
	}

	// strs should remain string
	strVals := result.Column("strs").Values()
	if v, ok := strVals[0].(string); !ok || v != "hello" {
		t.Errorf("ConvertDtypes: strs[0] = %v (%T), want hello", strVals[0], strVals[0])
	}
}
