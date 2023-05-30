package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gobackpack/rmq"
)

type Publisher struct {
	hub  *rmq.Hub
	conf map[string]*rmq.Publisher
}

func NewPublisher(ctx context.Context, hub *rmq.Hub, exchange string, events []string) *Publisher {
	pub := &Publisher{
		hub:  hub,
		conf: make(map[string]*rmq.Publisher),
	}

	pub.setupEvents(ctx, exchange, events)

	return pub
}

func (pub *Publisher) Publish(event string, msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	rmqPub := pub.conf[event]
	if rmqPub == nil {
		return errors.New("invalid event")
	}

	go rmqPub.Publish(b)

	return nil
}

func (pub *Publisher) setupEvents(ctx context.Context, exchange string, events []string) {
	for _, event := range events {
		conf := rmq.NewConfig()
		conf.Exchange = exchange
		conf.RoutingKey = event
		conf.ExchangeKind = "topic"

		pub.conf[event] = pub.hub.CreatePublisher(ctx, conf)
	}
}
