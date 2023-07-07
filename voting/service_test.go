package voting_test

import (
	"errors"
	"github.com/semirm-dev/mahala-backend/candidates"
	"github.com/semirm-dev/mahala-backend/datastore"
	"github.com/semirm-dev/mahala-backend/voting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterVotingTicket(t *testing.T) {
	dataStore := &datastore.MockDataStore{
		Candidates: []candidates.Candidate{
			{
				ID:   "candidate-1",
				Name: "candidate name",
			},
		},
	}

	err := voting.RegisterVotingTicket(dataStore, voting.Ticket{
		CandidateID: "candidate-1",
		VoterID:     "voter-1",
	}, &mockCandidatesApi{})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(dataStore.Votes))
	vote := dataStore.Votes[0]
	assert.Equal(t, "candidate-1", vote.CandidateID)
	assert.Equal(t, "voter-1", vote.VoterID)
}

func TestRegisterVotingTicket_CandidateNotExists(t *testing.T) {
	dataStore := &datastore.MockDataStore{}

	err := voting.RegisterVotingTicket(dataStore, voting.Ticket{
		CandidateID: "candidate-1",
		VoterID:     "voter-1",
	}, &mockCandidatesApi{withErr: errors.New("candidate not found")})
	assert.Equal(t, errors.New("candidate not found"), err)
	assert.Equal(t, 0, len(dataStore.Candidates))
	assert.Equal(t, 0, len(dataStore.Votes))
	assert.Equal(t, 0, len(dataStore.ProcessedVoters))
}

func TestQueryVotes(t *testing.T) {
	dataStore := &datastore.MockDataStore{
		Candidates: []candidates.Candidate{
			{
				ID:   "candidate-1",
				Name: "candidate name",
			},
		},
		Votes: []voting.Vote{
			{
				CandidateID: "candidate-1",
				VoterID:     "voter-1",
			},
		},
	}

	votes, err := voting.QueryVotes(dataStore, voting.QueryVoteFilter{CandidateID: "candidate-1"})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(votes))

	vote := votes[0]
	assert.Equal(t, "candidate-1", vote.CandidateID)
	assert.Equal(t, "voter-1", vote.VoterID)
}

type mockCandidatesApi struct {
	withErr error
}

func (api *mockCandidatesApi) GetCandidate(candidateID string) (voting.Candidate, error) {
	if api.withErr != nil {
		return voting.Candidate{}, api.withErr
	}
	return voting.Candidate{
		ID:   "candidate-1",
		Name: "candidate name",
	}, nil
}
