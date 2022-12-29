package main

import (
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/tachunwu/neural/pkg/jetstream"
)

func main() {
	// Connect NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	// Create streams
	jetstream.CreateStreams(nats.DefaultURL, Streams)

	// Calvin consumer Service A
	consumerA := NewCalvinConsumer(js, "CalvinConsumerA", "Calvin")
	consumerA.Start()

	// Calvin consumer Service B
	consumerB := NewCalvinConsumer(js, "CalvinConsumerB", "Calvin")
	consumerB.Start()

	// Run simulate client
	Client()

	for {
		runtime.Gosched()
	}
}
