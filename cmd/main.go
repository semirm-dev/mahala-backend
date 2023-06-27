package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/internal/pubsub"
	"github.com/semirm-dev/mahala/internal/redis"
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

	redisConf := redis.NewConfig()
	redisClient := redis.NewClient(redisConf)
	dataStore := voting.NewRedisStorage(redisClient)

	pubsub.Listen(hubCtx, hub, voting.HandleVotedEvent(dataStore))

	publisher := pubsub.NewPublisher(hubCtx, hub, voting.Bus, []string{voting.EventVoted})
	voteWriter := voting.PubSubVoteWriter(publisher)
	ticketSender := voting.NewTicketSender(voting.VoterValidator(dataStore), voteWriter)

	votesApi := api.Group("votes")
	votesApi.POST("", voting.SendVoteHandler(ticketSender))
	votesApi.GET("", voting.QueryVotesHandler(dataStore))

	candidatesApi := api.Group("candidates")
	candidatesApi.POST("", candidates.AddNewHandler(dataStore))
	candidatesApi.GET("", candidates.GetAllHandler(dataStore))

	web.ServeHttp(*httpAddr, "api", router)
}
