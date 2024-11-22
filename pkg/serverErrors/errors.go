package serverErrors

import "fmt"

var (
	INTERNAL_SERVER_ERROR = fmt.Errorf("The server encountered a problem and could not process your request")
)
