package usecase

import "fmt"

var (
	EMPTY_PASSWORD                = fmt.Errorf("password must be non-empty")
	INVALID_ROLE                  = fmt.Errorf("invalid role specified: must be either 'slave' or 'master'")
	EMPTY_SEARCH_STRING           = fmt.Errorf("an empty search string has been specified")
	INVALID_PRICE_RANGE           = fmt.Errorf("you have specified invalid price range: min and max prices must non-negative; min price must be less or equal to max price")
	EMPTY_TITLE                   = fmt.Errorf("empty title not allowed")
	NSFW_CONTENT_AVATAR_ERROR     = fmt.Errorf("avatar image you trying to publish seems to be an explicit content and not suitable for work")
	NSFW_CONTENT_BACK_IMAGE_ERROR = fmt.Errorf("back image you trying to publish seems to be an explicit content and not suitable for work")
	SWEAR_WORDS_ERROR             = fmt.Errorf("some of your input fileds contain insulting words")
	POSITIVE_NUMBER_REQUIRED      = fmt.Errorf("positive number required")
)
