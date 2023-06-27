package voting

import "errors"

type MockDataStore struct {
	Votes           []Vote
	ProcessedVoters []string
	Candidates      []string
}

func (ds *MockDataStore) StoreVote(candidateID string, votes []Vote) error {
	ds.Votes = append(ds.Votes, votes...)
	return nil
}

func (ds *MockDataStore) GetVotes(candidateID string) ([]Vote, error) {
	var votes []Vote

	for _, vote := range ds.Votes {
		if vote.CandidateID == candidateID {
			votes = append(votes, vote)
		}
	}

	return votes, nil
}

func (ds *MockDataStore) SetVoterAsProcessed(voterID string) error {
	ds.ProcessedVoters = append(ds.ProcessedVoters, voterID)
	return nil
}

func (ds *MockDataStore) GetProcessedVoters() ([]string, error) {
	return ds.ProcessedVoters, nil
}

func (ds *MockDataStore) AddCandidate(candidateID string) error {
	ds.Candidates = append(ds.Candidates, candidateID)
	return nil
}

func (ds *MockDataStore) GetCandidates() ([]string, error) {
	return ds.Candidates, nil
}

func (ds *MockDataStore) GetCandidate(candidateID string) (string, error) {
	for _, candidate := range ds.Candidates {
		if candidate == candidateID {
			return candidate, nil
		}
	}

	return "", errors.New("candidate not found")
}
