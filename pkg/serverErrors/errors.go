package serverErrors

import "fmt"

var (
	INTERNAL_SERVER_ERROR = fmt.Errorf("The server encountered a problem and could not process your request")
	CAST_ERROR            = fmt.Errorf("error while casting a variable to another type")

	SWEAR_WORDS_ERROR             = fmt.Errorf("some of your input fileds contain insulting words")
	NSFW_CONTENT_AVATAR_ERROR     = fmt.Errorf("avatar image you trying to publish seems to be an explicit content and not suitable for work")
	NSFW_CONTENT_BACK_IMAGE_ERROR = fmt.Errorf("back image you trying to publish seems to be an explicit content and not suitable for work")
)
