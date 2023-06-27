package candidates

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
	CandidateID string `json:"candidateID"`
}

func AddCandidateHandler(dataStore DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var candidateRequest NewCandidateRequest
		if err := c.ShouldBindJSON(&candidateRequest); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		if err := RegisterNew(dataStore, candidateRequest.CandidateID); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, HandlerResponse{Message: fmt.Sprintf("candidate %s created", candidateRequest.CandidateID)})
	}
}

func GetAllHandler(dataStore DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		candidates, err := GetAll(dataStore)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, candidates)
	}
}
