package coordinator

import (
	"github.com/code-sleuth/powerplant-monitoring-simulator/src/distributed/qutils"
	"github.com/streadway/amqp"
)

const url = "amqp://boss:inetutils@localhost:5672"

type QueueListener struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

func  NewQueueListener()  *QueueListener {
	ql := QueueListener{}

	ql.conn, ql. ch  = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("")
	_ = ql.ch.QueueBind(
		q.Name,       // name string,
		"",           // key string,
		"amq.fanout", // exchange string,
		false, // noWait bool,
		nil, // args amqp.Table
	)

	msgs, _ := ql.ch.Consume(
		q.Name, // queue string,
		"", // consumer sting,
		true, // utoAck bool,
		false, // exclusive bool,
		false, // noLocal bool,
		false, // noWait bool,
		nil, // args amqp.Table
	)

	for msg := range msgs {
		sourceChan, _ := ql.ch.Consume(
			string(msg.Body), // queue string,
			"", // consumer string,
			true, // utoAck bool,
			false, // exclusive bool,
			false, // noLocal bool,
			false, // noWait bool,
			nil, // args amqp.Table
		)
	}
}