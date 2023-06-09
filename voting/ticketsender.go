package voting

import (
	"errors"
	"github.com/semirm-dev/mahala-backend/internal/errwrapper"
	"github.com/sirupsen/logrus"
	"strings"
)

type TicketSender struct {
	validateVoter VoterValidatorFunc
	vote          VoteWriterFunc
}

// VoterValidatorFunc will validate voter (ID checks, has already voted...).
type VoterValidatorFunc func(voterID string) error

// VoteWriterFunc will send voting ticket to voting service
type VoteWriterFunc func(ticket Ticket) error

type Ticket struct {
	CandidateID string `json:"candidateID"`
	VoterID     string `json:"voterID"`
}

func NewTicketSender(voterIDValidator VoterValidatorFunc, voteWriter VoteWriterFunc) TicketSender {
	return TicketSender{
		validateVoter: voterIDValidator,
		vote:          voteWriter,
	}
}

// Send voting ticket to voting service
func (s TicketSender) Send(ticket Ticket) error {
	if err := isTicketValid(ticket); err != nil {
		return err
	}

	if err := s.validateVoter(ticket.VoterID); err != nil {
		return err
	}

	logrus.Infof("voter %s voting for %s", ticket.VoterID, ticket.CandidateID)

	return s.vote(ticket)
}

func isTicketValid(ticket Ticket) error {
	var err error

	if strings.TrimSpace(ticket.CandidateID) == "" {
		err = errwrapper.Wrap(err, errors.New("missing <candidateID>"))
	}

	if strings.TrimSpace(ticket.VoterID) == "" {
		err = errwrapper.Wrap(err, errors.New("missing <voterID>"))
	}

	return err
}
