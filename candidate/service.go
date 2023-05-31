package candidate

import "errors"

var (
	ErrCandidateExists = errors.New("candidate already registered")
)

// DataStore for candidates
type DataStore interface {
	AddCandidate(candidate string) error
	GetCandidates() ([]string, error)
}

// RegisterNew will add/create new candidate.
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

// GetAllCandidates currently registered.
func GetAllCandidates(dataStore DataStore) ([]string, error) {
	return dataStore.GetCandidates()
}
