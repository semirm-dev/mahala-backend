package pubsub

import (
	"context"
	"fmt"
	"github.com/gobackpack/rmq"
	"github.com/sirupsen/logrus"
)

type ListenerFunc func(ctx context.Context, hub *rmq.Hub)

func Listen(ctx context.Context, hub *rmq.Hub, listeners ...ListenerFunc) {
	for _, listenFunc := range listeners {
		go listenFunc(ctx, hub)
	}
}

func StartConsumer(ctx context.Context, hub *rmq.Hub, serviceName string, exchange, event string) *rmq.Consumer {
	conf := rmq.NewConfig()
	conf.Exchange = exchange
	conf.Queue = fmt.Sprintf("%s@%s", serviceName, event)
	conf.ConsumerTag = fmt.Sprintf("%s@%s", serviceName, event)
	conf.RoutingKey = event
	conf.ExchangeKind = "topic"

	if err := hub.CreateQueue(conf); err != nil {
		logrus.Fatal(err)
	}

	return hub.StartConsumer(ctx, conf)
}
