package candidates_test

import (
	"errors"
	"github.com/semirm-dev/mahala-backend/candidates"
	"github.com/semirm-dev/mahala-backend/datastore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterNew(t *testing.T) {
	dataStore := &datastore.MockDataStore{}

	err := candidates.RegisterNew(dataStore, candidates.Candidate{
		ID:           "candidate-1",
		Name:         "candidate name",
		ProfileImage: "img-1",
		Party:        "party-1",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dataStore.Candidates))

	candidate := dataStore.Candidates[0]
	assert.Equal(t, candidates.Candidate{
		ID:           "candidate-1",
		Name:         "candidate name",
		ProfileImage: "img-1",
		Party:        "party-1",
	}, candidate)
}

func TestRegisterNew_MissingCandidateID(t *testing.T) {
	dataStore := &datastore.MockDataStore{}

	err := candidates.RegisterNew(dataStore, candidates.Candidate{
		Name: "candidate name",
	})
	assert.Equal(t, errors.New("missing <candidate.ID>:missing <candidate.Party>"), err)
	assert.Equal(t, 0, len(dataStore.Candidates))
}

func TestRegisterNew_MissingCandidateName(t *testing.T) {
	dataStore := &datastore.MockDataStore{}

	err := candidates.RegisterNew(dataStore, candidates.Candidate{
		ID: "candidate-1",
	})
	assert.Equal(t, errors.New("missing <candidate.Name>:missing <candidate.Party>"), err)
	assert.Equal(t, 0, len(dataStore.Candidates))
}

func TestGetAll(t *testing.T) {
	dataStore := &datastore.MockDataStore{
		Candidates: []candidates.Candidate{
			{
				ID:           "candidate-1",
				Name:         "candidate name",
				ProfileImage: "img-1",
				Party:        "party-1",
			},
		},
	}

	allCandidates, err := candidates.GetAll(dataStore)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allCandidates))

	candidate := allCandidates[0]
	assert.Equal(t, "candidate-1", candidate.ID)
}
