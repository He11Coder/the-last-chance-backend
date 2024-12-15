package mongoTLC

import "fmt"

var (
	BAD_USER_ID           = fmt.Errorf("bad user ID")
	BAD_PET_ID            = fmt.Errorf("bad pet ID")
	BAD_SERVICE_ID        = fmt.Errorf("bad_service_id")
	NOT_FOUND             = fmt.Errorf("no data found")
	EMPTY_LOGIN           = fmt.Errorf("login must be non-empty")
	LOGIN_EXISTS          = fmt.Errorf("specified login already exists")
	INCORRECT_CREDENTIALS = fmt.Errorf("incorrect credentials")
	ACCESS_DENIED         = fmt.Errorf("you have no access to this resource")
)
