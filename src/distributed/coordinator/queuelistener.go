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
