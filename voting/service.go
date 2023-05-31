package voting

import (
	"context"
	"encoding/json"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/mahala/internal/pubsub"
	"github.com/sirupsen/logrus"
)

type Vote struct {
	Candidate string `json:"candidate"`
	VoterID   string `json:"voterID"`
}

// QueryVoteFilter is filter when querying votes.
type QueryVoteFilter struct {
	Candidate string `json:"candidate"`
}

// DataStore is used to store votes
type DataStore interface {
	Store(candidate string, votes []Vote) error
	Get(candidate string) ([]Vote, error)
}

// VotedEventHandler Reacts to EventVoted event. It will handle all votes requests, store votes in datastore.
func VotedEventHandler(dataStore DataStore) func(ctx context.Context, hub *rmq.Hub) {
	return func(ctx context.Context, hub *rmq.Hub) {
		defer logrus.Warnf("%s consumer closed", EventVoted)

		errors := make(chan error)
		go handleErrors(ctx, errors)

		consumer := pubsub.StartConsumer(ctx, hub, "voting_service", Bus, string(EventVoted))

		logrus.Infof("%s started", EventVoted)

		for {
			select {
			case msg := <-consumer.OnMessage:
				var ticket Ticket
				if err := json.Unmarshal(msg, &ticket); err != nil {
					errors <- err
				}

				votes, err := dataStore.Get(ticket.VoteFor)
				if err != nil {
					errors <- err
					continue
				}

				votes = append(votes, Vote{
					Candidate: ticket.VoteFor,
					VoterID:   ticket.VoterID,
				})

				if err := dataStore.Store(ticket.VoteFor, votes); err != nil {
					errors <- err
					continue
				}

				logrus.Infof("[%s] - user %s successfully voted for %s", EventVoted, ticket.VoterID, ticket.VoteFor)
			case err := <-consumer.OnError:
				errors <- err
			case <-ctx.Done():
				return
			}
		}
	}
}

// QueryVotes returns either all or filtered votes.
func QueryVotes(filter QueryVoteFilter) ([]Vote, error) {
	return nil, nil
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
