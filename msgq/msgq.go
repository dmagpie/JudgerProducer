package msgq

import (
	"JudgerProducer/config"

	"github.com/streadway/amqp"
)

// Conn RMQ连接
var Conn *amqp.Connection

func init() {
	amqpuri, err := config.GetConfig("RMQ_URL")
	if err != nil {
		panic(err)
	}
	c, err := amqp.Dial(amqpuri)
	if err != nil {
		panic(err)
	}
	Conn = c
}
