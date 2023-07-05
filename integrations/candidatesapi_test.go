package integrations_test

import (
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/datastore"
	"github.com/semirm-dev/mahala/integrations"
	"github.com/semirm-dev/mahala/voting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCandidatesApiAdapter_GetCandidate(t *testing.T) {
	candidatesAdapter := integrations.CandidatesApiAdapter{
		CandidatesDataStore: &datastore.MockDataStore{
			Candidates: []candidates.Candidate{
				{
					ID:   "candidate-1",
					Name: "candidate name",
				},
			},
		},
	}

	votingCandidate, err := candidatesAdapter.GetCandidate("candidate-1")
	assert.NoError(t, err)

	expectedCandidate := voting.Candidate{
		ID:   "candidate-1",
		Name: "candidate name",
	}
	assert.Equal(t, expectedCandidate, votingCandidate)
}

func TestCandidatesApiAdapter_GetCandidate_NotExists(t *testing.T) {
	candidatesAdapter := integrations.CandidatesApiAdapter{
		CandidatesDataStore: &datastore.MockDataStore{
			Candidates: []candidates.Candidate{
				{
					ID:   "candidate-1",
					Name: "candidate name",
				},
			},
		},
	}

	votingCandidate, err := candidatesAdapter.GetCandidate("candidate-2")
	assert.Error(t, err)

	expectedCandidate := voting.Candidate{}
	assert.Equal(t, expectedCandidate, votingCandidate)
}
