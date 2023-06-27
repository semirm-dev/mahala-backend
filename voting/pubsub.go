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

// PubSubVoteWriter will publish an event of EventVoted for new voting ticket.
func PubSubVoteWriter(pub *pubsub.Publisher) VoteWriterFunc {
	return func(ticket Ticket) error {
		return pub.Publish(EventVoted, ticket)
	}
}

// HandleVotedEvent reacts to EventVoted event. Calls voting business logic.
func HandleVotedEvent(dataStore DataStore) func(ctx context.Context, hub *rmq.Hub) {
	return func(ctx context.Context, hub *rmq.Hub) {
		defer logrus.Warnf("consumer for event %s closed", EventVoted)

		errors := make(chan error)
		go handleErrors(ctx, errors)

		consumer := pubsub.StartConsumer(ctx, hub, "voting_service", Bus, string(EventVoted))

		logrus.Infof("consumer for event %s started", EventVoted)

		for {
			select {
			case msg := <-consumer.OnMessage:
				var ticket Ticket
				if err := json.Unmarshal(msg, &ticket); err != nil {
					errors <- err
				}

				if err := RegisterVotingTicket(dataStore, ticket); err != nil {
					errors <- err
				}
			case err := <-consumer.OnError:
				errors <- err
			case <-ctx.Done():
				return
			}
		}
	}
}

func handleErrors(ctx context.Context, errors chan error) {
	for {
		select {
		case <-ctx.Done():
			break
		case err := <-errors:
			logrus.Error(err)
		}
	}
}
