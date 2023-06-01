package voting

type Vote struct {
	Candidate string `json:"candidate"`
	VoterID   string `json:"voterID"`
}

// QueryVoteFilter is filter when querying Votes.
type QueryVoteFilter struct {
	Candidate string `json:"candidate"`
}

// DataStore is used to store Votes
type DataStore interface {
	StoreVote(candidate string, votes []Vote) error
	GetVotes(candidate string) ([]Vote, error)
	SetVoterAsProcessed(voterID string) error
	GetProcessedVoters() ([]string, error)
}

func RegisterVotingTicket(dataStore DataStore, ticket Ticket) error {
	if err := vote(dataStore, ticket); err != nil {
		return err
	}

	if err := setVoterAsProcessed(dataStore, ticket.VoterID); err != nil {
		return err
	}

	return nil
}

// QueryVotes returns either all or filtered Votes.
func QueryVotes(dataStore DataStore, filter QueryVoteFilter) ([]Vote, error) {
	votes, err := dataStore.GetVotes(filter.Candidate)
	if err != nil {
		return nil, err
	}

	return votes, nil
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
