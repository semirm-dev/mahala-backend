package voting

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HandlerResponse struct {
	Message string `json:"message"`
}

type QueryVotesResponse struct {
	Votes []Vote `json:"Votes"`
}

func SendVoteHandler(ticketSender TicketSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket Ticket
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

		c.JSON(http.StatusOK, HandlerResponse{Message: "successfully voted"})
	}
}

func QueryVotesHandler(dataStore DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter QueryVoteFilter
		if err := c.ShouldBindJSON(&filter); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		votes, err := QueryVotes(dataStore, filter)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, QueryVotesResponse{
			Votes: votes,
		})
	}
}
