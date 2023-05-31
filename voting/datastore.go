package voting

import (
	"encoding/json"
	"github.com/semirm-dev/mahala/internal/redis"
)

const (
	processedVotersKey = "processed_voters"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) RedisStorage {
	return RedisStorage{
		redisClient: redisClient,
	}
}

func (r RedisStorage) StoreVote(candidate string, votes []Vote) error {
	return r.redisClient.Add(redis.Item{
		Key:   candidate,
		Value: votes,
	})
}

func (r RedisStorage) GetVotes(candidate string) ([]Vote, error) {
	candidateVotes, err := r.redisClient.Get(candidate)
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
	processedVoters, err := r.redisClient.Get(processedVotersKey)
	if err != nil && err != redis.ErrNotExists {
		return err
	}

	var processed []string
	if len(processedVoters) > 0 {
		if err := json.Unmarshal(processedVoters, &processed); err != nil {
			return err
		}
	}

	processed = append(processed, voterID)

	return r.redisClient.Add(redis.Item{
		Key:   processedVotersKey,
		Value: processed,
	})
}

func (r RedisStorage) GetProcessedVoters() ([]string, error) {
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
