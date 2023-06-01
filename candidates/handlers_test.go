package candidates_test

import (
	"bytes"
	"encoding/json"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/internal/web"
	"github.com/semirm-dev/mahala/voting"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCandidateHandler(t *testing.T) {
	dataStore := &voting.MockDataStore{}
	router := web.NewRouter()
	router.POST("/", candidates.AddCandidateHandler(dataStore))

	payload := `{"candidate": "candidate-1"}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := candidates.HandlerResponse{Message: "candidate candidate-1 created"}
	var addCandidateResponse candidates.HandlerResponse

	err := json.NewDecoder(w.Body).Decode(&addCandidateResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, addCandidateResponse)
}

func TestGetAllHandler(t *testing.T) {
	dataStore := &voting.MockDataStore{
		Candidates: []string{"candidate-1", "candidate-2"},
	}
	router := web.NewRouter()
	router.GET("/", candidates.GetAllHandler(dataStore))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	expectedResponse := []string{"candidate-1", "candidate-2"}
	var candidatesResponse []string

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
