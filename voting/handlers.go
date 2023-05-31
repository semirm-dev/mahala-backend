package voting

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type VoteResponse struct {
	Successful bool `json:"successful"`
}

type QueryVotesResponse struct {
	Votes []Vote `json:"votes"`
}

func VoteHandler(ticketSender TicketSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := ticketSender.Send(ticket); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, VoteResponse{Successful: true})
	}
}

func QueryVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter QueryVoteFilter
		if err := c.ShouldBindJSON(&filter); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		votes, err := QueryVotes(filter)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, QueryVotesResponse{
			Votes: votes,
		})
	}
}
