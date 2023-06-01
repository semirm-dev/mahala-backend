package voting_test

import (
	"bytes"
	"encoding/json"
	"github.com/semirm-dev/mahala/internal/web"
	"github.com/semirm-dev/mahala/voting"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVoteHandler(t *testing.T) {
	router := web.NewRouter()
	ticketSender := voting.NewTicketSender(fakeVoterIDValidator, fakeVoteWriter)
	router.POST("/", voting.VoteHandler(ticketSender))

	payload := `{"voterID": "voter-123", "voteFor": "candidate-123"}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := voting.HandlerResponse{Message: "vote successful"}
	var voteResponse voting.HandlerResponse

	err := json.NewDecoder(w.Body).Decode(&voteResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, voteResponse)
}

func TestQueryVoteHandler(t *testing.T) {
	dataStore := &voting.MockDataStore{
		Votes: []voting.Vote{
			{
				Candidate: "candidate-1",
				VoterID:   "voter-123",
			},
		},
	}

	payload := `{"candidate": "candidate-1"}`

	router := web.NewRouter()
	router.GET("/", voting.QueryVoteHandler(dataStore))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := voting.QueryVotesResponse{
		Votes: []voting.Vote{
			{
				Candidate: "candidate-1",
				VoterID:   "voter-123",
			},
		},
	}
	var queryResponse voting.QueryVotesResponse

	err := json.NewDecoder(w.Body).Decode(&queryResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, queryResponse)
}

func mockHttpServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"message": "success"}
		res, _ := json.Marshal(response)
		_, err := w.Write(res)
		assert.Nil(t, err)
	}))
}
