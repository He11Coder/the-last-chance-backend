package mongoTLC

import "fmt"

var (
	BAD_USER_ID = fmt.Errorf("bad user ID")
	NOT_FOUND   = fmt.Errorf("no data found")
)
