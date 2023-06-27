package voting

import (
	"encoding/json"
	"github.com/semirm-dev/mahala/internal/redis"
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

func (r RedisStorage) StoreVote(candidateID string, votes []Vote) error {
	return r.redisClient.Add(redis.Item{
		Key:   candidateID,
		Value: votes,
	})
}

func (r RedisStorage) GetVotes(candidateID string) ([]Vote, error) {
	candidateVotes, err := r.redisClient.Get(candidateID)
	if err != nil && err != redis.ErrNotExists {
		return nil, err
	}

	var votes []Vote
	if len(candidateVotes) > 0 {
		if err := json.Unmarshal(candidateVotes, &votes); err != nil {
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

func (r RedisStorage) AddCandidate(candidateID string) error {
	candidates, err := r.getCandidates()
	if err != nil {
		return err
	}

	candidates = append(candidates, candidateID)

	return r.redisClient.Add(redis.Item{
		Key:   candidatesKey,
		Value: candidates,
	})
}

func (r RedisStorage) GetCandidates() ([]string, error) {
	return r.getCandidates()
}

func (r RedisStorage) getProcessedVoters() ([]string, error) {
	processedVoters, err := r.redisClient.Get(processedVotersKey)
	if err != nil && err != redis.ErrNotExists {
		return nil, err
	}

	var processed []string
	if len(processedVoters) > 0 {
		if err := json.Unmarshal(processedVoters, &processed); err != nil {
			return nil, err
		}
	}

	return processed, nil
}

func (r RedisStorage) getCandidates() ([]string, error) {
	existingCandidates, err := r.redisClient.Get(candidatesKey)
	if err != nil && err != redis.ErrNotExists {
		return nil, err
	}

	var candidates []string
	if len(existingCandidates) > 0 {
		if err := json.Unmarshal(existingCandidates, &candidates); err != nil {
			return nil, err
		}
	}
	return candidates, nil
}
