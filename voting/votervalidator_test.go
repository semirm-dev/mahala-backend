package voting_test

import (
	"errors"
	"github.com/semirm-dev/mahala/voting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoterValidator(t *testing.T) {
	testTable := map[string]struct {
		processedVoters []string
		expectedErr     error
	}{
		"not processed voter returns nil error": {
			processedVoters: nil,
			expectedErr:     nil,
		},
		"already processed voter returns error": {
			processedVoters: []string{"voter-1"},
			expectedErr:     errors.New("voter voter-1 has already voted"),
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			dataStore := &mockDataStore{
				processedVoters: tt.processedVoters,
			}

			voterValidator := voting.VoterValidator(dataStore)

			err := voterValidator("voter-1")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
