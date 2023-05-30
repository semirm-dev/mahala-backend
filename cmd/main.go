package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/semirm-dev/mahala/internal/web"
	"github.com/semirm-dev/mahala/voting"
	"net/http"
)

var (
	httpAddr = flag.String("http", ":8000", "Http address")
)

func main() {
	flag.Parse()

	router := web.NewRouter()

	api := router.Group("api")

	api.GET("healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})

	votes := api.Group("votes")

	ticketSender := voting.NewTicketSender(voting.FakeVoterIDValidator, voting.PubSubVoteWriter)

	votes.POST("", voting.VoteHandler(ticketSender))

	web.ServeHttp(*httpAddr, "api", router)
}
