package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/code-sleuth/powerplant-monitoring-simulator/src/distributed/dto"
	"github.com/code-sleuth/powerplant-monitoring-simulator/src/distributed/qutils"
	"github.com/streadway/amqp"
)
var url = os.Getenv("CONNSTRING") //"amqp://<user>:<password>@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
}

func NewQueueListener() *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch)
	_ = ql.ch.QueueBind(
		q.Name,       // name string,
		"",           // key string,
		"amq.fanout", // exchange string,
		false,        // noWait bool,
		nil,          // args amqp.Table
	)

	msgs, _ := ql.ch.Consume(
		q.Name, // queue string,
		"",     // consumer sting,
		true,   // utoAck bool,
		false,  // exclusive bool,
		false,  // noLocal bool,
		false,  // noWait bool,
		nil,    // args amqp.Table
	)

	for msg := range msgs {
		sourceChan, _ := ql.ch.Consume(
			string(msg.Body), // queue string,
			"",               // consumer string,
			true,             // utoAck bool,
			false,            // exclusive bool,
			false,            // noLocal bool,
			false,            // noWait bool,
			nil,              // args amqp.Table
		)

		if ql.sources[string(msg.Body)] == nil {
			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Recieved message: %v\n", sd)
	}
}
