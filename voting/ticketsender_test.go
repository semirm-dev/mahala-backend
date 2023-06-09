package voting_test

import (
	"errors"
	"fmt"
	"github.com/semirm-dev/mahala-backend/voting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTicketSender_Send(t *testing.T) {
	testTable := map[string]struct {
		ticket           voting.Ticket
		voterIdValidator voting.VoterValidatorFunc
		voteWriter       voting.VoteWriterFunc
		expectedErr      error
	}{
		"ticket successfully sent": {
			ticket: voting.Ticket{
				CandidateID: "candidate-1",
				VoterID:     "voter-123",
			},
			voterIdValidator: fakeVoterIDValidator,
			voteWriter:       fakeVoteWriter,
			expectedErr:      nil,
		},
		"missing candidateID": {
			ticket: voting.Ticket{
				VoterID: "voter-123",
			},
			voterIdValidator: fakeVoterIDValidator,
			voteWriter:       fakeVoteWriter,
			expectedErr:      errors.New("missing <candidateID>"),
		},
		"missing voterID": {
			ticket: voting.Ticket{
				CandidateID: "candidate-1",
			},
			voterIdValidator: fakeVoterIDValidator,
			voteWriter:       fakeVoteWriter,
			expectedErr:      errors.New("missing <voterID>"),
		},
		"voter id validator with error should return error": {
			ticket: voting.Ticket{
				CandidateID: "candidate-1",
				VoterID:     "voter-123",
			},
			voterIdValidator: func(voterID string) error {
				return errors.New(fmt.Sprintf("voter id is invalid"))
			},
			voteWriter:  fakeVoteWriter,
			expectedErr: errors.New(fmt.Sprintf("voter id is invalid")),
		},
		"applyVote writer with error should return error": {
			ticket: voting.Ticket{
				CandidateID: "candidate-1",
				VoterID:     "voter-123",
			},
			voterIdValidator: fakeVoterIDValidator,
			voteWriter: func(ticket voting.Ticket) error {
				return errors.New("applyVote writer failed")
			},
			expectedErr: errors.New("applyVote writer failed"),
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			ticketSender := voting.NewTicketSender(tt.voterIdValidator, tt.voteWriter)

			err := ticketSender.Send(tt.ticket)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func fakeVoteWriter(ticket voting.Ticket) error {
	return nil
}

func fakeVoterIDValidator(voterID string) error {
	return nil
}
