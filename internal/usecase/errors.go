package usecase

import "fmt"

var (
	EMPTY_PASSWORD                = fmt.Errorf("password must be non-empty")
	INVALID_ROLE                  = fmt.Errorf("invalid role specified: must be either 'slave' or 'master'")
	EMPTY_SEARCH_STRING           = fmt.Errorf("an empty search string has been specified")
	EMPTY_TITLE                   = fmt.Errorf("empty title not allowed")
	NSFW_CONTENT_AVATAR_ERROR     = fmt.Errorf("avatar image you trying to publish seems to be explicit and not suitable for work")
	NSFW_CONTENT_BACK_IMAGE_ERROR = fmt.Errorf("back image you trying to publish seems to be explicit and not suitable for work")
)
