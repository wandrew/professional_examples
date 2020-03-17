package amqp

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
)

func (q *Queue) Deliver(ch *rabbitmq.Channel) <-chan amqp.Delivery {
	queue, err := ch.QueueDeclare(
		q.Name, // name of the queue
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // noWait
		nil,    // arguments
	)
	if err != nil {
		panic(fmt.Errorf("Queue Declare: %s", err))
	}

	if err = ch.QueueBind(
		queue.Name,      // name of the queue
		q.BindingKey,    // bindingKey
		q.Exchange.Name, // sourceExchange
		false,           // noWait
		nil,             // arguments
	); err != nil {
		panic(fmt.Errorf("Queue Bind: %s", err))

	}

	deliveries, err := ch.Consume(
		queue.Name,    // name
		q.ConsumerTag, // consumerTag,
		true,          // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)

	if err != nil {
		panic(fmt.Errorf("Delivery: %s", err))
	}

	return deliveries
}

func (e Exchange) Declare(ch *rabbitmq.Channel) (*rabbitmq.Channel, error) {
	if err := ch.ExchangeDeclare(
		e.Name, // name of the exchange
		e.Type, // type
		true,   // durable
		false,  // delete when complete
		false,  // service
		false,  // noWait
		nil,    // arguments
	); err != nil {
		return nil, err
	}

	return ch, nil
}
