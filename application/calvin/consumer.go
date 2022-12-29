package main

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/tachunwu/neural/pkg/jetstream"
)

type CalvinConsumer struct {
	ctx       context.Context
	jetstream nats.JetStreamContext
	durable   string
	subject   string
}

func NewCalvinConsumer(
	js nats.JetStreamContext,
	durable string,
	subject string,
) *CalvinConsumer {

	ctx, _ := context.WithCancel(context.Background())

	return &CalvinConsumer{
		ctx:       ctx,
		jetstream: js,
		durable:   durable,
		subject:   subject,
	}
}

func (c *CalvinConsumer) Start() error {
	return jetstream.JetStreamConsumer(
		c.ctx,
		c.jetstream,
		c.subject,
		c.durable,
		1,
		c.HandleFn,
		nats.DeliverAll(),
		nats.ManualAck(),
	)
}

func (c *CalvinConsumer) HandleFn(ctx context.Context, msgs []*nats.Msg) bool {
	for _, msg := range msgs {

		log.Println("Service:", c.durable, "local commit:", msg.Header.Get("Nats-Msg-Id"))

		commit := &nats.Msg{
			Subject: "Calvin." + msg.Header.Get("Nats-Msg-Id"),
			Header:  nats.Header{},
		}
		commit.Header.Add("Service-Id", c.durable)
		c.jetstream.PublishMsg(commit)
	}
	return true
}
