package voting

import (
	"errors"
	"fmt"
)

func VoterIDValidator(dataStore DataStore) VoterIDValidatorFunc {
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
