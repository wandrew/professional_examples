package amqp

import (
	"github.com/streadway/amqp"
)

type Exchange struct {
	Name string
	Type string
}

type Queue struct {
	Name        string
	BindingKey  string
	ConsumerTag string
	Exchange    Exchange
}

type ServiceBroker struct {
	Channel    *amqp.Channel
	Exchange   Exchange
	Queue      Queue
	Deliveries chan<- amqp.Delivery
}

// This needs to be added back into service package Shouldn't be in here.

// type SubscriberProfile struct {
// 	Logger   log.Logger
// 	Endpoint endpoint.Endpoints
// 	Channel  *amqp.Channel
// 	Exchange Exchange
// 	Queue    Queue
// }
