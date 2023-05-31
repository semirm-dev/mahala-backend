package voting

import (
	"errors"
	"fmt"
)

func VoterValidator(dataStore DataStore) VoterValidatorFunc {
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
