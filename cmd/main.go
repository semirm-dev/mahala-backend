package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gobackpack/rmq"
	candidatesApi "github.com/semirm-dev/mahala-backend/candidates/api"
	"github.com/semirm-dev/mahala-backend/datastore"
	"github.com/semirm-dev/mahala-backend/integrations"
	"github.com/semirm-dev/mahala-backend/internal/pubsub"
	"github.com/semirm-dev/mahala-backend/internal/redis"
	"github.com/semirm-dev/mahala-backend/internal/web"
	"github.com/semirm-dev/mahala-backend/voting"
	votingApi "github.com/semirm-dev/mahala-backend/voting/api"
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
	dataStore := datastore.NewRedisStorage(redisClient)

	pubsub.Listen(hubCtx, hub, integrations.HandleVotedEvent(dataStore, integrations.CandidatesApiAdapter{CandidatesDataStore: dataStore}))

	publisher := pubsub.NewPublisher(hubCtx, hub, integrations.Bus, []string{integrations.EventVoted})
	voteWriter := integrations.PubSubVoteWriter(publisher)
	ticketSender := voting.NewTicketSender(integrations.VoterValidator(dataStore), voteWriter)

	votes := api.Group("votes")
	votes.POST("", votingApi.SendVoteHandler(ticketSender))
	votes.GET("", votingApi.QueryVotesHandler(dataStore))

	candidates := api.Group("candidates")
	candidates.POST("", candidatesApi.AddNewCandidateHandler(dataStore))
	candidates.GET("", candidatesApi.GetAllCandidatesHandler(dataStore))

	web.ServeHttp(*httpAddr, "api", router)
}
