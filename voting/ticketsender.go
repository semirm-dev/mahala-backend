package voting

import "github.com/sirupsen/logrus"

type TicketSender struct {
	validateVoterID VoterIDValidatorFunc
	vote            VoteWriterFunc
}

// VoterIDValidatorFunc will validate voter (ID...).
type VoterIDValidatorFunc func(voterID string) error

// VoteWriterFunc will send voting ticket to voting service
type VoteWriterFunc func(ticket Ticket) error

type Ticket struct {
	VoteFor string `json:"voteFor"`
	VoterID string `json:"voterID"`
}

func NewTicketSender(voterIDValidator VoterIDValidatorFunc, voteWriter VoteWriterFunc) TicketSender {
	return TicketSender{
		validateVoterID: voterIDValidator,
		vote:            voteWriter,
	}
}

// Send voting ticket to voting service
func (s TicketSender) Send(ticket Ticket) error {
	if err := s.validateVoterID(ticket.VoterID); err != nil {
		return err
	}

	logrus.Infof("voter %s voting for %s", ticket.VoterID, ticket.VoteFor)

	return s.vote(ticket)
}
