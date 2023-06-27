package candidates

import "errors"

var (
	ErrCandidateExists = errors.New("candidate already registered")
)

// DataStore for candidates
type DataStore interface {
	AddCandidate(candidateID string) error
	GetCandidates() ([]string, error)
}

// RegisterNew new candidate.
func RegisterNew(dataStore DataStore, candidateID string) error {
	existingCandidates, err := dataStore.GetCandidates()
	if err != nil {
		return err
	}

	for _, existingCandidate := range existingCandidates {
		if existingCandidate == candidateID {
			return ErrCandidateExists
		}
	}

	return dataStore.AddCandidate(candidateID)
}

// GetAll currently registered candidates.
func GetAll(dataStore DataStore) ([]string, error) {
	return dataStore.GetCandidates()
}
