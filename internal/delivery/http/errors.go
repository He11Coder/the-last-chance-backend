package http

import "fmt"

var (
	BAD_QUERY_PARAMETERS = fmt.Errorf("Bad Query parameters specified")
	BAD_GET_PARAMETER    = fmt.Errorf("Bad GET parameter specified")
	INVALID_BODY         = fmt.Errorf("Bad request body")
	BAD_JSON_FORMAT      = fmt.Errorf("invalid json format: must be with fields 'username' and 'password'")
	MISSING_USER_ID      = fmt.Errorf("user ID is missing")
	AUTH_ERROR           = fmt.Errorf("authorization error")
	AVATAR_ERROR         = fmt.Errorf("error while reading user's avatar")
)
