package jetstream

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func JetStreamConsumer(
	ctx context.Context,
	js nats.JetStreamContext,
	subject string,
	durable string,
	batch int,
	consumeFn func(ctx context.Context, msgs []*nats.Msg) bool,
	opts ...nats.SubOpt,
) error {
	defer func() {}()

	sub, err := js.PullSubscribe(subject, durable, opts...)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			default:
			}

			msgs, err := sub.Fetch(batch, nats.Context(ctx))

			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Second)
			}

			if len(msgs) < 1 {
				continue
			}

			for _, msg := range msgs {
				err := msg.InProgress(nats.Context(ctx))
				if err != nil {
					log.Println(err)
					continue
				}
			}

			if consumeFn(ctx, msgs) {
				for _, msg := range msgs {
					err := msg.AckSync(nats.Context(ctx))
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				for _, msg := range msgs {
					err := msg.Nak(nats.Context(ctx))
					if err != nil {
						log.Println(err)
					}
				}
			}

		}
	}()

	return nil
}

func CreateStreams(addr string, streams []*nats.StreamConfig) {
	// Connect to NATS
	nc, _ := nats.Connect(addr)
	defer nc.Drain()
	js, _ := nc.JetStream()

	for _, stream := range streams {
		_, err := js.AddStream(stream)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
