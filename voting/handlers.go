package voting

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type VoteResponse struct {
	Successful bool `json:"successful"`
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
