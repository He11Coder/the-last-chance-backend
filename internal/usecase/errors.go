package usecase

import "fmt"

var (
	EMPTY_PASSWORD = fmt.Errorf("password must be non-empty")
	INVALID_ROLE   = fmt.Errorf("invalid role specified: must be either 'slave' or 'master'")
)
