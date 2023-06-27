package candidates

import (
	"errors"
	"github.com/semirm-dev/mahala/internal/errwrapper"
	"strings"
)

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
	if err := isCandidateValid(candidateID); err != nil {
		return err
	}

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

func isCandidateValid(candidateID string) error {
	var err error

	if strings.TrimSpace(candidateID) == "" {
		err = errwrapper.Wrap(err, errors.New("missing <candidateID>"))
	}

	return err
}
