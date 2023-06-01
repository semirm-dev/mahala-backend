package voting

type MockDataStore struct {
	Votes           []Vote
	ProcessedVoters []string
	Candidates      []string
}

func (ds *MockDataStore) StoreVote(candidate string, votes []Vote) error {
	ds.Votes = append(ds.Votes, votes...)
	return nil
}

func (ds *MockDataStore) GetVotes(candidate string) ([]Vote, error) {
	var votes []Vote

	for _, vote := range ds.Votes {
		if vote.Candidate == candidate {
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

func (ds *MockDataStore) AddCandidate(candidate string) error {
	ds.Candidates = append(ds.Candidates, candidate)
	return nil
}

func (ds *MockDataStore) GetCandidates() ([]string, error) {
	return ds.Candidates, nil
}
