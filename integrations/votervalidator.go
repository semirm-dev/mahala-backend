package integrations

import (
	"errors"
	"fmt"
	"github.com/semirm-dev/mahala/voting"
)

func VoterValidator(dataStore voting.DataStore) voting.VoterValidatorFunc {
	return func(voterID string) error {
		voted, err := hasVoted(dataStore, voterID)
		if err != nil {
			return err
		}

		if voted {
			return errors.New(fmt.Sprintf("voter %s has already voted", voterID))
		}

		return nil
	}
}

func hasVoted(dataStore voting.DataStore, voterID string) (bool, error) {
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
