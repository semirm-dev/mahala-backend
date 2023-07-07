package integrations

import (
	"github.com/semirm-dev/mahala-backend/candidates"
	"github.com/semirm-dev/mahala-backend/voting"
)

type CandidatesApiAdapter struct {
	CandidatesDataStore candidates.DataStore
}

func (api CandidatesApiAdapter) GetCandidate(candidateID string) (voting.Candidate, error) {
	candidate, err := candidates.GetByID(api.CandidatesDataStore, candidateID)
	if err != nil {
		return voting.Candidate{}, err
	}

	return voting.Candidate{
		ID:   candidate.ID,
		Name: candidate.Name,
	}, nil
}
