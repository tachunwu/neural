package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

func Client() {

	// Connect NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	// Publish message
	n := nuid.New()

	nc.Subscribe("Calvin.*", func(msg *nats.Msg) {
		log.Println("Coordinator commit:", msg.Header.Get("Service-Id"))
	})

	id := n.Next()
	msg := &nats.Msg{
		Subject: "Calvin",
		Header:  nats.Header{},
	}
	msg.Header.Add("Nats-Msg-Id", id)
	// Subscribe own response

	// Publish transaction
	js.PublishMsg(msg)

}
