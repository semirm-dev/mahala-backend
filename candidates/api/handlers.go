package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HandlerResponse struct {
	Message string `json:"message"`
}

type NewCandidateRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProfileImage string `json:"profileImage"`
	Party        string `json:"party"`
}

func AddNewCandidateHandler(dataStore candidates.DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var candidateRequest NewCandidateRequest
		if err := c.ShouldBindJSON(&candidateRequest); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		if err := candidates.RegisterNew(dataStore, candidates.Candidate{
			ID:           candidateRequest.ID,
			Name:         candidateRequest.Name,
			ProfileImage: candidateRequest.ProfileImage,
			Party:        candidateRequest.Party,
		}); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, HandlerResponse{Message: fmt.Sprintf("candidate %s created", candidateRequest.Name)})
	}
}

func GetAllCandidatesHandler(dataStore candidates.DataStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		candidatesResponse, err := candidates.GetAll(dataStore)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, HandlerResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, candidatesResponse)
	}
}
