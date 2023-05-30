package voting

import (
	"context"
	"encoding/json"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/mahala/internal/pubsub"
	"github.com/sirupsen/logrus"
)

const (
	Bus        = "voting_bus"
	EventVoted = "event_voted"
)

func PubSubVoteWriter(pub *pubsub.Publisher) VoteWriterFunc {
	return func(ticket Ticket) error {
		return pub.Publish(EventVoted, ticket)
	}
}

func VotedEventHandler(ctx context.Context, hub *rmq.Hub) {
	defer logrus.Warnf("%s consumer closed", EventVoted)

	consumer := pubsub.StartConsumer(ctx, hub, "voting_service", Bus, string(EventVoted))

	logrus.Infof("%s started", EventVoted)

	for {
		select {
		case msg := <-consumer.OnMessage:
			var ticket Ticket
			json.Unmarshal(msg, &ticket)

			logrus.Infof("[%s] - user %s successfully voted for %s", EventVoted, ticket.VoterID, ticket.VoteFor)
		case err := <-consumer.OnError:
			logrus.Error(err)
		case <-ctx.Done():
			return
		}
	}
}
