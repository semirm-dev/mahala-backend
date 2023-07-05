package api_test

import (
	"bytes"
	"encoding/json"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/candidates/api"
	"github.com/semirm-dev/mahala/datastore"
	"github.com/semirm-dev/mahala/internal/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCandidateHandler(t *testing.T) {
	dataStore := &datastore.MockDataStore{}
	router := web.NewRouter()
	router.POST("/", api.AddNewCandidateHandler(dataStore))

	payload := `{"id": "candidate-1", "name": "candidate name"}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := api.HandlerResponse{Message: "candidate candidate name created"}
	var addCandidateResponse api.HandlerResponse

	err := json.NewDecoder(w.Body).Decode(&addCandidateResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, addCandidateResponse)
}

func TestGetAllCandidatesHandler(t *testing.T) {
	dataStore := &datastore.MockDataStore{
		Candidates: []candidates.Candidate{
			{
				ID:   "candidate-1",
				Name: "candidate name",
			},
			{
				ID:   "candidate-2",
				Name: "candidate 2 name",
			},
		},
	}
	router := web.NewRouter()
	router.GET("/", api.GetAllCandidatesHandler(dataStore))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := []candidates.Candidate{
		{
			ID:   "candidate-1",
			Name: "candidate name",
		},
		{
			ID:   "candidate-2",
			Name: "candidate 2 name",
		},
	}
	var candidatesResponse []candidates.Candidate

	err := json.NewDecoder(w.Body).Decode(&candidatesResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, candidatesResponse)
}

func mockHttpServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"message": "success"}
		res, _ := json.Marshal(response)
		_, err := w.Write(res)
		assert.Nil(t, err)
	}))
}
