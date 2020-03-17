package amqp

import (
	"fmt"

	"github.com/streadway/amqp"

	"github.com/deseretdigital/sl-colony-platform/config"
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
)

var Broker struct {
	ConnectionString string
}

func init() {
	config.ConfigureService(&Broker, "amqp")
}

func Connect() (*rabbitmq.Connection, error) {

	c, err := rabbitmq.Dial(Broker.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.NotifyClose(make(chan *amqp.Error)))
	}()

	//	TODO: More pointer weirdness
	return c, nil
}
