package mongoTLC

import "fmt"

var (
	BAD_USER_ID = fmt.Errorf("bad user ID")
	BAD_PET_ID  = fmt.Errorf("bad pet ID")
	NOT_FOUND   = fmt.Errorf("no data found")
)
