package readwrite

import (
	"fmt"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ReadParquet reads a DataFrame from a Parquet file.
// Parquet binary format support is not yet implemented. Use CSV instead.
func ReadParquet(filename string) (*tabgo.DataFrame, error) {
	return nil, fmt.Errorf("readwrite: Parquet read not implemented: " +
		"Parquet binary format not yet supported; use CSV instead")
}

// WriteParquet writes a DataFrame to a Parquet file.
// Parquet binary format support is not yet implemented. Use CSV instead.
func WriteParquet(filename string, df *tabgo.DataFrame) error {
	return fmt.Errorf("readwrite: Parquet write not implemented: " +
		"Parquet binary format not yet supported; use CSV instead")
}
