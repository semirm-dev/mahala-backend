package datastore

import (
	"encoding/json"
	"errors"
	"github.com/semirm-dev/mahala/candidates"
	"github.com/semirm-dev/mahala/internal/redis"
	"github.com/semirm-dev/mahala/voting"
)

const (
	processedVotersKey = "processed_voters"
	candidatesKey      = "candidates"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) RedisStorage {
	return RedisStorage{
		redisClient: redisClient,
	}
}

func (r RedisStorage) StoreVote(candidateID string, votes []voting.Vote) error {
	return r.redisClient.Add(redis.Item{
		Key:   candidateID,
		Value: votes,
	})
}

func (r RedisStorage) GetVotes(candidateID string) ([]voting.Vote, error) {
	candidateVotes, err := r.redisClient.Get(candidateID)
	if err != nil && err != redis.ErrNotExists {
		return nil, err
	}

	var votes []voting.Vote
	if len(candidateVotes) > 0 {
		if err = json.Unmarshal(candidateVotes, &votes); err != nil {
			return nil, err
		}
	}

	return votes, nil
}

func (r RedisStorage) SetVoterAsProcessed(voterID string) error {
	processed, err := r.getProcessedVoters()
	if err != nil {
		return err
	}

	processed = append(processed, voterID)

	return r.redisClient.Add(redis.Item{
		Key:   processedVotersKey,
		Value: processed,
	})
}

func (r RedisStorage) GetProcessedVoters() ([]string, error) {
	return r.getProcessedVoters()
}

func (r RedisStorage) AddCandidate(candidate candidates.Candidate) error {
	existingCandidates, err := r.getCandidates()
	if err != nil {
		return err
	}

	existingCandidates = append(existingCandidates, candidate)

	return r.redisClient.Add(redis.Item{
		Key:   candidatesKey,
		Value: existingCandidates,
	})
}

func (r RedisStorage) GetCandidates() ([]candidates.Candidate, error) {
	return r.getCandidates()
}

func (r RedisStorage) GetCandidate(candidateID string) (candidates.Candidate, error) {
	var candidate candidates.Candidate

	unmarshalledCandidates, err := r.redisClient.Get(candidatesKey)
	if err != nil || err == redis.ErrNotExists || len(unmarshalledCandidates) == 0 {
		return candidate, errors.New("candidate not found")
	}

	var existingCandidates []candidates.Candidate
	if err = json.Unmarshal(unmarshalledCandidates, &existingCandidates); err != nil {
		return candidate, err
	}

	for _, c := range existingCandidates {
		if c.ID == candidateID {
			candidate = c
		}
	}

	return candidate, nil
}

func (r RedisStorage) getProcessedVoters() ([]string, error) {
	processedVoters, err := r.redisClient.Get(processedVotersKey)
	if err != nil && err != redis.ErrNotExists {
		return nil, err
	}

	var processed []string
	if len(processedVoters) > 0 {
		if err = json.Unmarshal(processedVoters, &processed); err != nil {
			return nil, err
		}
	}

	return processed, nil
}

func (r RedisStorage) getCandidates() ([]candidates.Candidate, error) {
	existingCandidates, err := r.redisClient.Get(candidatesKey)
	if err != nil && err != redis.ErrNotExists {
		return nil, err
	}

	var candidatesResult []candidates.Candidate
	if len(existingCandidates) > 0 {
		if err := json.Unmarshal(existingCandidates, &candidatesResult); err != nil {
			return nil, err
		}
	}
	return candidatesResult, nil
}
