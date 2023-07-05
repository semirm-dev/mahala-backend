package voting

import (
	"errors"
)

type Vote struct {
	CandidateID string `json:"candidateID"`
	VoterID     string `json:"voterID"`
}

// QueryVoteFilter is filter when querying Votes.
type QueryVoteFilter struct {
	CandidateID string `json:"candidateID"`
}

type Candidate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DataStore is used to store Votes
type DataStore interface {
	StoreVote(candidateID string, votes []Vote) error
	GetVotes(candidateID string) ([]Vote, error)
	SetVoterAsProcessed(voterID string) error
	GetProcessedVoters() ([]string, error)
}

// CandidatesApi is responsible for communication with candidates service/domain.
type CandidatesApi interface {
	GetCandidate(candidateID string) (Candidate, error)
}

func RegisterVotingTicket(dataStore DataStore, ticket Ticket, candidatesApi CandidatesApi) error {
	candidate, err := candidatesApi.GetCandidate(ticket.CandidateID)
	if err != nil {
		return err
	}

	if candidate.ID == "" || candidate.Name == "" {
		return errors.New("candidate not found")
	}

	if err = applyVote(dataStore, ticket); err != nil {
		return err
	}

	if err = setVoterAsProcessed(dataStore, ticket.VoterID); err != nil {
		return err
	}

	return nil
}

// QueryVotes returns either all or filtered Votes.
func QueryVotes(dataStore DataStore, filter QueryVoteFilter) ([]Vote, error) {
	votes, err := dataStore.GetVotes(filter.CandidateID)
	if err != nil {
		return nil, err
	}

	return votes, nil
}

func applyVote(dataStore DataStore, ticket Ticket) error {
	votes, err := dataStore.GetVotes(ticket.CandidateID)
	if err != nil {
		return err
	}

	votes = append(votes, Vote{
		CandidateID: ticket.CandidateID,
		VoterID:     ticket.VoterID,
	})

	return dataStore.StoreVote(ticket.CandidateID, votes)
}

func setVoterAsProcessed(dataStore DataStore, voterID string) error {
	return dataStore.SetVoterAsProcessed(voterID)
}
