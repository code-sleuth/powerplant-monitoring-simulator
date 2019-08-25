package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"
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
)


func main() {
	flag.Parse()

	duration, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(duration)

	for range signal {
		calcValue()
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
