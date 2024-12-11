package redisTLC

import "fmt"

var (
	SESSION_NOT_FOUND = fmt.Errorf("no session corresponding to this ID was found")
)
