/*
Gocelery is a task queue implementation for Go modules used to asynchronously execute work outside the HTTP request-response cycle. Celery is an implementation of the task queue concept.
*/

package gocelery

import (
	"log"

	// "github.com/streadway/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	FirstErrorNotice bool = true
)

// deliveryAck acknowledges delivery message with retries on error
func deliveryAck(delivery amqp.Delivery) {
	var err error
	for retryCount := 3; retryCount > 0; retryCount-- {
		if err = delivery.Ack(false); err == nil {
			break
		}
	}

	if err != nil && FirstErrorNotice {
		FirstErrorNotice = false
		log.Printf("The connection suddenly closed between the rabbitMQ!")
	}
}
