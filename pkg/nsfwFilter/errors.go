package nsfwFilter

import "fmt"

var (
	IMAGE_MARSHAL_ERR         = fmt.Errorf("failed to marshal image into json")
	RABBIT_CONNECT_ERR        = fmt.Errorf("failed to connect to RabbitMQ")
	CHANNEL_OPENNING_ERR      = fmt.Errorf("failed to open a channel")
	QUEUE_DECLARATION_ERR     = fmt.Errorf("failed to declare a queue")
	CONSUMER_REGISTRATION_ERR = fmt.Errorf("failed to register a consumer")
	MESSAGE_PUBLISHING_ERR    = fmt.Errorf("failed to publish a message")
	RESPONSE_CONVERTION_ERR   = fmt.Errorf("failed to convert response body to json")

	NO_RESPONSE_CODE        = fmt.Errorf("no response code has been returned")
	ERR_CASTING_CODE        = fmt.Errorf("error while casting response code to type int")
	NO_ERR_DETAILS          = fmt.Errorf("expected error description has not been received")
	ERR_CASTING_ERR_DETAILS = fmt.Errorf("error while casting error details to type string")
	NO_SAFE_FLAG            = fmt.Errorf("no boolean response (safe/not safe for work) has been returned")
	ERR_CASTING_SAFE_FLAG   = fmt.Errorf("error while casting boolean response (safe/not safe for work) to type bool")
	NO_CONFIDENCE           = fmt.Errorf("expected confidence level has not been received")
	ERR_CASTING_CONFIDENCE  = fmt.Errorf("error while casting confidence level to type float64")

	REQUEST_ERROR = fmt.Errorf("a bad request has been done, for more information see Inference.Err field")
	WORKER_ERROR  = fmt.Errorf("nsfw validator error happened, for more information see Inference.Err field")

	TIMEOUT_ERR = fmt.Errorf("a response has not been received before the timeout")

	BAD_REQUEST           = 400
	INTERNAL_SERVER_ERROR = 500
)
