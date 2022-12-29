package main

import "github.com/nats-io/nats.go"

var Streams = []*nats.StreamConfig{
	{
		Name:      "HAwriter",
		Retention: nats.InterestPolicy,
		Storage:   nats.FileStorage,
	},
}
