package serverErrors

import "fmt"

var (
	INTERNAL_SERVER_ERROR = fmt.Errorf("The server encountered a problem and could not process your request")
	CAST_ERROR            = fmt.Errorf("error while casting a variable to another type")
)
