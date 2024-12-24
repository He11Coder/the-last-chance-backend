package nsfwFilter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func errWithDetails(baseErr, details error) error {
	return fmt.Errorf("err: %v\n in detail: %v", baseErr, details)
}

type rawInference struct {
	IsSafe     bool    `json:"is_safe,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Code       int     `json:"code"`
	StrErr     string  `json:"error,omitempty"`
}

type Inference struct {
	IsSafe     bool
	Confidence float64
	Code       int
	Err        error
}

func parseResults(rawRes rawInference) (Inference, error) {
	result := Inference{}

	if rawRes.Code == 0 {
		return result, NO_RESPONSE_CODE
	}
	result.Code = rawRes.Code

	if rawRes.Code == BAD_REQUEST || rawRes.Code == INTERNAL_SERVER_ERROR {
		if rawRes.StrErr == "" {
			return result, NO_ERR_DETAILS
		}
		result.Err = fmt.Errorf("%s", rawRes.StrErr)

		if rawRes.Code == BAD_REQUEST {
			return result, REQUEST_ERROR
		}

		if rawRes.Code == INTERNAL_SERVER_ERROR {
			return result, WORKER_ERROR
		}
	}

	result.IsSafe = rawRes.IsSafe
	result.Confidence = rawRes.Confidence

	return result, nil
}

func IsSafeForWork(base64Image string) (Inference, error) {
	imageMap := map[string]interface{}{
		"image": base64Image,
	}

	jsonImage, err := json.Marshal(imageMap)
	if err != nil {
		return Inference{}, errWithDetails(IMAGE_MARSHAL_ERR, err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:8001/")
	if err != nil {
		return Inference{}, errWithDetails(RABBIT_CONNECT_ERR, err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return Inference{}, errWithDetails(CHANNEL_OPENNING_ERR, err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return Inference{}, errWithDetails(QUEUE_DECLARATION_ERR, err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return Inference{}, errWithDetails(CONSUMER_REGISTRATION_ERR, err)
	}

	corrId := uuid.NewString()

	err = ch.Publish(
		"",                      // exchange
		"nsfw_validation_queue", // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			DeliveryMode:  amqp.Persistent,
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          jsonImage,
		},
	)
	if err != nil {
		return Inference{}, errWithDetails(MESSAGE_PUBLISHING_ERR, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := rawInference{}
	for {
		select {
		case d := <-msgs:
			if corrId == d.CorrelationId {
				err = json.Unmarshal(d.Body, &result)
				if err != nil {
					return Inference{}, errWithDetails(RESPONSE_CONVERTION_ERR, err)
				}
			}
			return parseResults(result)
		case <-ctx.Done():
			return Inference{}, TIMEOUT_ERR
		}
	}
}
