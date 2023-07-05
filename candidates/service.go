package candidates

import (
	"errors"
	"github.com/semirm-dev/mahala/internal/errwrapper"
	"strings"
)

var (
	ErrCandidateExists = errors.New("candidate already registered")
)

type Candidate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DataStore is used to store Votes
type DataStore interface {
	AddCandidate(candidate Candidate) error
	GetCandidates() ([]Candidate, error)
	GetCandidate(candidateID string) (Candidate, error)
}

// RegisterNew new candidate.
func RegisterNew(dataStore DataStore, candidate Candidate) error {
	if err := isCandidateValid(candidate); err != nil {
		return err
	}

	existingCandidates, err := dataStore.GetCandidates()
	if err != nil {
		return err
	}

	for _, existingCandidate := range existingCandidates {
		if existingCandidate.ID == candidate.ID {
			return ErrCandidateExists
		}
	}

	return dataStore.AddCandidate(candidate)
}

// GetAll currently registered candidates.
func GetAll(dataStore DataStore) ([]Candidate, error) {
	return dataStore.GetCandidates()
}

func GetByID(dataStore DataStore, candidateID string) (Candidate, error) {
	return dataStore.GetCandidate(candidateID)
}

func isCandidateValid(candidate Candidate) error {
	var err error

	if strings.TrimSpace(candidate.ID) == "" {
		err = errwrapper.Wrap(err, errors.New("missing <candidate.ID>"))
	}

	if strings.TrimSpace(candidate.Name) == "" {
		err = errwrapper.Wrap(err, errors.New("missing <candidate.Name>"))
	}

	return err
}
