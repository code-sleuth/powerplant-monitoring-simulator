package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/code-sleuth/powerplant-monitoring-simulator/src/distributed/dto"
	"github.com/code-sleuth/powerplant-monitoring-simulator/src/distributed/qutils"
)

var (
	name = flag.String("name", "sensor", "name of the sensor")
	freq = flag.Uint("freq", 5, "update frequency in cycles per second")
	max = flag.Float64("max", 5., "maximum value for generated readings")
	min = flag.Float64("min", 1., "minimum value for generated readings")
	stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	value = r.Float64() * (*max - *min) + *min
	nom = (*max - *min) / 2 + *min

	url = os.Getenv("CONNSTRING") //"amqp://<user>:<password>@localhost:5672"

)


func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)
	// create queue to ensure messages are recieved
	//sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)



	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"amq.fanout", // exchange string,
		"", // key string,
		false, // mandatory bool,
		false, // immediate bool,
		msg, // msg amqp.Publishing,
		)

	duration, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(duration)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name: *name,
			Value: value,
			Timestamp: time.Now(),

		}
		buf.Reset()
		enc = gob.NewEncoder(buf)
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"", // exchange string,
			dataQueue.Name, // key string,
			false, // mandatory bool,
			false, // immediate bool,
			msg, // msg amqp.Publishing,
			)

		log.Printf("Reading sent. Value: %v\n", value)
	}
}

func calcValue(){
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64() * (maxStep - minStep) + minStep
}
