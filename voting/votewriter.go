package voting

import "github.com/semirm-dev/mahala/internal/pubsub"

const (
	Bus        = "voting_bus"
	EventVoted = "event_voted"
)

func PubSubVoteWriter(pub *pubsub.Publisher) VoteWriterFunc {
	return func(ticket Ticket) error {
		return pub.Publish(EventVoted, ticket)
	}
}
