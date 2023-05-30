package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/mahala/internal/pubsub"
	"github.com/semirm-dev/mahala/internal/web"
	"github.com/semirm-dev/mahala/voting"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	httpAddr = flag.String("http", ":8000", "Http address")
	rmqHost  = flag.String("rmq_host", "localhost", "RabbitMQ host address")
)

func main() {
	flag.Parse()

	router := web.NewRouter()

	api := router.Group("api")

	api.GET("healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})

	cred := rmq.NewCredentials()
	cred.Host = *rmqHost
	hub := rmq.NewHub(cred)

	hubCtx, hubCancel := context.WithCancel(context.Background())
	defer hubCancel()

	_, err := hub.Connect(hubCtx)
	if err != nil {
		logrus.Fatal(err)
	}

	pubsub.Listen(hubCtx, hub, voting.VotedEventHandler)

	publisher := pubsub.NewPublisher(hubCtx, hub, voting.Bus, []string{voting.EventVoted})
	voteWriter := voting.PubSubVoteWriter(publisher)

	ticketSender := voting.NewTicketSender(voting.VoterIDValidator, voteWriter)

	votes := api.Group("votes")
	votes.POST("", voting.VoteHandler(ticketSender))

	web.ServeHttp(*httpAddr, "api", router)
}
