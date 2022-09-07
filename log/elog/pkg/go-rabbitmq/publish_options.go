package gorabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishOptFunc func(*PublishOptions)

type PublishOptions struct {
	Exchange        string
	Mandatory       bool
	Immediate       bool
	ContentType     string
	DeliveryMode    uint8
	Expiration      string
	ContentEncoding string
	Priority        uint8
	CorrelationID   string
	ReplyTo         string
	MessageID       string
	Timestamp       time.Time
	Type            string
	UserID          string
	AppID           string
	Headers         amqp.Table
}
