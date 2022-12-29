package main

import "github.com/nats-io/nats.go"

var Streams = []*nats.StreamConfig{
	{
		Name:      "Calvin",
		Retention: nats.InterestPolicy,
		Storage:   nats.FileStorage,
	},
}
