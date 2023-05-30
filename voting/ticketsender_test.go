package voting_test

import (
	"errors"
	"fmt"
	"github.com/semirm-dev/mahala/voting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTicketSender_Send(t *testing.T) {
	testTable := map[string]struct {
		voterIdValidator voting.VoterIDValidator
		voteWriter       voting.VoteWriter
		expectedErr      error
	}{
		"ticket successfully sent": {
			voterIdValidator: voting.FakeVoterIDValidator,
			voteWriter:       voting.FakeVoteWriter,
			expectedErr:      nil,
		},
		"voter id validator with error should return error": {
			voterIdValidator: func(voterID string) error {
				return errors.New(fmt.Sprintf("voter id is invalid"))
			},
			voteWriter:  voting.FakeVoteWriter,
			expectedErr: errors.New(fmt.Sprintf("voter id is invalid")),
		},
		"vote writer with error should return error": {
			voterIdValidator: voting.FakeVoterIDValidator,
			voteWriter: func(ticket voting.Ticket) error {
				return errors.New("vote writer failed")
			},
			expectedErr: errors.New("vote writer failed"),
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			ticketSender := voting.NewTicketSender(tt.voterIdValidator, tt.voteWriter)

			err := ticketSender.Send(voting.Ticket{})
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
