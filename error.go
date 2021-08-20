package safer

import "errors"

var (
	// ErrCompanyNotFound is thrown when a company is not found for the searched MC/MX/DOT number
	ErrCompanyNotFound = errors.New("company not found")
)
