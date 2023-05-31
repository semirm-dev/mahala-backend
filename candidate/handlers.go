package candidate

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HandlerResponse struct {
	Message string `json:"message"`
}

type NewCandidateRequest struct {
	Candidate string `json:"candidate"`
}

func AddCandidateHandler(dataStore DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var candidateRequest NewCandidateRequest
		if err := c.ShouldBindJSON(&candidateRequest); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		if err := RegisterNew(dataStore, candidateRequest.Candidate); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, HandlerResponse{Message: fmt.Sprintf("candidate %s created", candidateRequest.Candidate)})
	}
}

func GetCandidatesHandler(dataStore DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		candidates, err := GetAllCandidates(dataStore)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, candidates)
	}
}
