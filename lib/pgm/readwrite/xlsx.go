package readwrite

import (
	"fmt"

	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

// ReadXLSX reads a BayesianNetwork from an XLSX file.
// XLSX support requires third-party ZIP/XML parsing libraries and is not yet
// implemented. Use JSON or CSV formats instead.
func ReadXLSX(filename string) (*models.BayesianNetwork, error) {
	return nil, fmt.Errorf("readwrite: XLSX read not implemented: " +
		"XLSX requires third-party ZIP/XML parsing; use JSON or CSV instead")
}

// WriteXLSX writes a BayesianNetwork to an XLSX file.
// XLSX support requires third-party ZIP/XML parsing libraries and is not yet
// implemented. Use JSON or CSV formats instead.
func WriteXLSX(filename string, bn *models.BayesianNetwork) error {
	return fmt.Errorf("readwrite: XLSX write not implemented: " +
		"XLSX requires third-party ZIP/XML parsing; use JSON or CSV instead")
}
