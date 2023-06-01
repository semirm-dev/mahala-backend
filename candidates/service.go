package candidates

import "errors"

var (
	ErrCandidateExists = errors.New("candidate already registered")
)

// DataStore for candidates
type DataStore interface {
	AddCandidate(candidate string) error
	GetCandidates() ([]string, error)
}

// RegisterNew new candidate.
func RegisterNew(dataStore DataStore, candidate string) error {
	existingCandidates, err := dataStore.GetCandidates()
	if err != nil {
		return err
	}

	for _, existingCandidate := range existingCandidates {
		if existingCandidate == candidate {
			return ErrCandidateExists
		}
	}

	return dataStore.AddCandidate(candidate)
}

// GetAll currently registered candidates.
func GetAll(dataStore DataStore) ([]string, error) {
	return dataStore.GetCandidates()
}
