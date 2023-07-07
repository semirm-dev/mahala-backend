package api

import (
	"github.com/gin-gonic/gin"
	"github.com/semirm-dev/mahala-backend/voting"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HandlerResponse struct {
	Message string `json:"message"`
}

func SendVoteHandler(ticketSender voting.TicketSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket voting.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		if err := ticketSender.Send(ticket); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, HandlerResponse{Message: "vote finished, will be evaluated"})
	}
}

func QueryVotesHandler(dataStore voting.DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		candidateID := c.Query("candidateID")

		votes, err := voting.QueryVotes(dataStore, voting.QueryVoteFilter{
			CandidateID: candidateID,
		})
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, votes)
	}
}
