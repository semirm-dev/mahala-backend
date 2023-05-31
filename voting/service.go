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
	StoreVote(candidate string, votes []Vote) error
	GetVotes(candidate string) ([]Vote, error)
	SetVoterAsProcessed(voterID string) error
	GetProcessedVoters() ([]string, error)
}

// HandleVotedEvent reacts to EventVoted event. It will handle all votes requests, store votes in datastore.
func HandleVotedEvent(dataStore DataStore) func(ctx context.Context, hub *rmq.Hub) {
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

				if err := vote(dataStore, ticket); err != nil {
					errors <- err
					continue
				}

				logrus.Infof("[%s] - user %s successfully voted for %s", EventVoted, ticket.VoterID, ticket.VoteFor)

				if err := setVoterAsProcessed(dataStore, ticket.VoterID); err != nil {
					errors <- err
					continue
				}
			case err := <-consumer.OnError:
				errors <- err
			case <-ctx.Done():
				return
			}
		}
	}
}

// QueryVotes returns either all or filtered votes.
func QueryVotes(dataStore DataStore, filter QueryVoteFilter) ([]Vote, error) {
	votes, err := dataStore.GetVotes(filter.Candidate)
	if err != nil {
		return nil, err
	}

	return votes, nil
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

func vote(dataStore DataStore, ticket Ticket) error {
	votes, err := dataStore.GetVotes(ticket.VoteFor)
	if err != nil {
		return err
	}

	votes = append(votes, Vote{
		Candidate: ticket.VoteFor,
		VoterID:   ticket.VoterID,
	})

	return dataStore.StoreVote(ticket.VoteFor, votes)
}

func setVoterAsProcessed(dataStore DataStore, voterID string) error {
	return dataStore.SetVoterAsProcessed(voterID)
}

func hasVoted(dataStore DataStore, voterID string) (bool, error) {
	voters, err := dataStore.GetProcessedVoters()
	if err != nil {
		return false, err
	}

	for _, voter := range voters {
		if voter == voterID {
			return true, nil
		}
	}

	return false, nil
}
