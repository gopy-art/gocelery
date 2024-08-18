// Copyright (c) 2019 Sick Yoon
// This file is part of gocelery which is released under MIT license.
// See file LICENSE for full license details.

package gocelery

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"zsploit/logger"

	// "github.com/streadway/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPExchange stores AMQP Exchange configuration
type AMQPExchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
}

// NewAMQPExchange creates new AMQPExchange
func NewAMQPExchange(name string) *AMQPExchange {
	return &AMQPExchange{
		Name:       name,
		Type:       "direct",
		Durable:    true,
		AutoDelete: true,
	}
}

// AMQPQueue stores AMQP Queue configuration
type AMQPQueue struct {
	Name       string
	Durable    bool
	AutoDelete bool
}

// NewAMQPQueue creates new AMQPQueue
func NewAMQPQueue(name string, durable bool) *AMQPQueue {
	return &AMQPQueue{
		Name:       name,
		Durable:    durable,
		AutoDelete: false,
	}
}

// AMQPCeleryBroker is RedisBroker for AMQP
type AMQPCeleryBroker struct {
	*amqp.Channel
	Connection       *amqp.Connection
	Exchange         *AMQPExchange
	Queue            *AMQPQueue
	consumingChannel <-chan amqp.Delivery
	Rate             int
}

// NewAMQPConnection creates new AMQP channel
func NewAMQPConnection(host string) (*amqp.Connection, *amqp.Channel) {
	var rabbitMQConnenctionStatusCeleryErrorMessage bool = true
	for {
		connection, err := amqp.Dial(host)
		if err != nil {
			if rabbitMQConnenctionStatusCeleryErrorMessage {
				logger.ErrorLogger.Printf("we have an error to connecting to this url : %v!\nThe error is = %v", host, err)
				rabbitMQConnenctionStatusCeleryErrorMessage = !rabbitMQConnenctionStatusCeleryErrorMessage
			}
		} else {
			rabbitMQConnenctionStatusCeleryErrorMessage = !rabbitMQConnenctionStatusCeleryErrorMessage
			channel, err := connection.Channel()
			if err != nil {
				logger.ErrorLogger.Printf("we have an error to connecting to this url : %v!\nThe error is = %v", host, err)
			}
			_ = rabbitMQConnenctionStatusCeleryErrorMessage
			
			return connection, channel
		}

		secondNumber, _ := strconv.Atoi(os.Getenv("RABBIT_INTERVAL"))
		time.Sleep(time.Second * time.Duration(secondNumber))
	}
}

// NewAMQPCeleryBroker creates new AMQPCeleryBroker
func NewAMQPCeleryBroker(host string, queue string, durable bool) *AMQPCeleryBroker {
	connection, channel := NewAMQPConnection(host)
	return NewAMQPCeleryBrokerByConnAndChannel(connection, channel, queue, durable)
}

// NewAMQPCeleryBrokerByConnAndChannel creates new AMQPCeleryBroker using AMQP conn and channel
func NewAMQPCeleryBrokerByConnAndChannel(conn *amqp.Connection, channel *amqp.Channel, queue string, durable bool) *AMQPCeleryBroker {
	broker := &AMQPCeleryBroker{
		Channel:    channel,
		Connection: conn,
		Exchange:   NewAMQPExchange("default"),
		Queue:      NewAMQPQueue(queue, durable),
		Rate:       1,
	}
	if err := broker.CreateExchange(); err != nil {
		panic(err)
	}
	if err := broker.CreateQueue(); err != nil {
		panic(err)
	}
	if err := broker.Qos(broker.Rate, 0, false); err != nil {
		panic(err)
	}
	if err := broker.StartConsumingChannel(); err != nil {
		panic(err)
	}
	return broker
}

// StartConsumingChannel spawns receiving channel on AMQP queue
func (b *AMQPCeleryBroker) StartConsumingChannel() error {
	channel, err := b.Consume(b.Queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	b.consumingChannel = channel
	return nil
}

// SendCeleryMessage sends CeleryMessage to broker
func (b *AMQPCeleryBroker) SendCeleryMessage(message *CeleryMessage) error {
	taskMessage := message.GetTaskMessage()
	queueName := b.Queue.Name
	_, err := b.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return err
	}
	err = b.ExchangeDeclare(
		"default",
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	resBytes, err := json.Marshal(taskMessage)
	if err != nil {
		return err
	}

	publishMessage := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         resBytes,
	}

	return b.Publish(
		"",
		queueName,
		false,
		false,
		publishMessage,
	)
}

// GetTaskMessage retrieves task message from AMQP queue
func (b *AMQPCeleryBroker) GetTaskMessage() (*TaskMessage, error) {
	select {
	case delivery := <-b.consumingChannel:
		deliveryAck(delivery)
		var taskMessage TaskMessage
		if err := json.Unmarshal(delivery.Body, &taskMessage); err != nil {
			return nil, err
		}
		return &taskMessage, nil
	default:
		return nil, fmt.Errorf("consuming channel is empty")
	}
}

// CreateExchange declares AMQP exchange with stored configuration
func (b *AMQPCeleryBroker) CreateExchange() error {
	return b.ExchangeDeclare(
		b.Exchange.Name,
		b.Exchange.Type,
		b.Exchange.Durable,
		b.Exchange.AutoDelete,
		false,
		false,
		nil,
	)
}

// CreateQueue declares AMQP Queue with stored configuration
func (b *AMQPCeleryBroker) CreateQueue() error {
	_, err := b.QueueDeclare(
		b.Queue.Name,
		b.Queue.Durable,
		b.Queue.AutoDelete,
		false,
		false,
		nil,
	)
	return err
}
