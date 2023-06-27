package candidates_test

import (
	"errors"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/voting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterNew(t *testing.T) {
	dataStore := &voting.MockDataStore{}

	err := candidates.RegisterNew(dataStore, "candidate-1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dataStore.Candidates))

	candidate := dataStore.Candidates[0]
	assert.Equal(t, "candidate-1", candidate)
}

func TestRegisterNew_MissingCandidateID(t *testing.T) {
	dataStore := &voting.MockDataStore{}

	err := candidates.RegisterNew(dataStore, "")
	assert.Equal(t, errors.New("missing <candidateID>"), err)
	assert.Equal(t, 0, len(dataStore.Candidates))
}

func TestGetAll(t *testing.T) {
	dataStore := &voting.MockDataStore{
		Candidates: []string{"candidate-1"},
	}

	allCandidates, err := candidates.GetAll(dataStore)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allCandidates))

	candidate := allCandidates[0]
	assert.Equal(t, "candidate-1", candidate)
}
