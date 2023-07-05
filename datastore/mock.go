package datastore

import (
	"errors"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/voting"
)

type MockDataStore struct {
	Votes           []voting.Vote
	ProcessedVoters []string
	Candidates      []candidates.Candidate
}

func (ds *MockDataStore) StoreVote(candidateID string, votes []voting.Vote) error {
	ds.Votes = append(ds.Votes, votes...)
	return nil
}

func (ds *MockDataStore) GetVotes(candidateID string) ([]voting.Vote, error) {
	var votes []voting.Vote

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

func (ds *MockDataStore) AddCandidate(candidate candidates.Candidate) error {
	ds.Candidates = append(ds.Candidates, candidate)
	return nil
}

func (ds *MockDataStore) GetCandidates() ([]candidates.Candidate, error) {
	return ds.Candidates, nil
}

func (ds *MockDataStore) GetCandidate(candidateID string) (candidates.Candidate, error) {
	for _, candidate := range ds.Candidates {
		if candidate.ID == candidateID {
			return candidate, nil
		}
	}

	return candidates.Candidate{}, errors.New("candidate not found")
}
