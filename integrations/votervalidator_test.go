package integrations_test

import (
	"errors"
	"github.com/semirm-dev/mahala-backend/datastore"
	"github.com/semirm-dev/mahala-backend/integrations"
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
			dataStore := &datastore.MockDataStore{
				ProcessedVoters: tt.processedVoters,
			}

			voterValidator := integrations.VoterValidator(dataStore)

			err := voterValidator("voter-1")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
