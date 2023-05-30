package voting

import "github.com/sirupsen/logrus"

type TicketSender struct {
	validateVoterID VoterIDValidator
	vote            VoteWriter
}

// VoterIDValidator will validate voter's ID.
type VoterIDValidator func(voterID string) error

// VoteWriter will send voting ticket to voting service
type VoteWriter func(ticket Ticket) error

type Ticket struct {
	VoterID  string `json:"voterID"`
	VotedFor string `json:"votedFor"`
}

func NewTicketSender(voterIDValidator VoterIDValidator, voteWriter VoteWriter) TicketSender {
	return TicketSender{
		validateVoterID: voterIDValidator,
		vote:            voteWriter,
	}
}

func (s TicketSender) Send(ticket Ticket) error {
	if err := s.validateVoterID(ticket.VoterID); err != nil {
		return err
	}

	logrus.Infof("voter %s voting for %s", ticket.VoterID, ticket.VotedFor)

	return s.vote(ticket)
}
