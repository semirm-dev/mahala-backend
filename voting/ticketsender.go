package voting

import "github.com/sirupsen/logrus"

type TicketSender struct {
	validateVoter VoterValidatorFunc
	vote          VoteWriterFunc
}

// VoterValidatorFunc will validate voter (ID checks, has already voted...).
type VoterValidatorFunc func(voterID string) error

// VoteWriterFunc will send voting ticket to voting service
type VoteWriterFunc func(ticket Ticket) error

type Ticket struct {
	VoteFor string `json:"voteFor"`
	VoterID string `json:"voterID"`
}

func NewTicketSender(voterIDValidator VoterValidatorFunc, voteWriter VoteWriterFunc) TicketSender {
	return TicketSender{
		validateVoter: voterIDValidator,
		vote:          voteWriter,
	}
}

// Send voting ticket to voting service
func (s TicketSender) Send(ticket Ticket) error {
	if err := s.validateVoter(ticket.VoterID); err != nil {
		return err
	}

	logrus.Infof("voter %s voting for %s", ticket.VoterID, ticket.VoteFor)

	return s.vote(ticket)
}
