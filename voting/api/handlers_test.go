package api_test

import (
	"bytes"
	"encoding/json"
	"github.com/semirm-dev/mahala-backend/datastore"
	"github.com/semirm-dev/mahala-backend/internal/web"
	"github.com/semirm-dev/mahala-backend/voting"
	votingApi "github.com/semirm-dev/mahala-backend/voting/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVoteHandler(t *testing.T) {
	router := web.NewRouter()
	ticketSender := voting.NewTicketSender(fakeVoterIDValidator, fakeVoteWriter)
	router.POST("/", votingApi.SendVoteHandler(ticketSender))

	payload := `{"voterID": "voter-123", "candidateID": "candidate-123"}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := votingApi.HandlerResponse{Message: "vote finished, will be evaluated"}
	var voteResponse votingApi.HandlerResponse

	err := json.NewDecoder(w.Body).Decode(&voteResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, voteResponse)
}

func TestQueryVoteHandler(t *testing.T) {
	dataStore := &datastore.MockDataStore{
		Votes: []voting.Vote{
			{
				CandidateID: "candidate-1",
				VoterID:     "voter-123",
			},
		},
	}

	router := web.NewRouter()
	router.GET("/", votingApi.QueryVotesHandler(dataStore))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/?candidateID=candidate-1", nil)
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := []voting.Vote{
		{
			CandidateID: "candidate-1",
			VoterID:     "voter-123",
		},
	}
	var queryResponse []voting.Vote

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

func fakeVoterIDValidator(voterID string) error {
	return nil
}

func fakeVoteWriter(ticket voting.Ticket) error {
	return nil
}
