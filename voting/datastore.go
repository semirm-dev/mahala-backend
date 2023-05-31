package voting

import (
	"encoding/json"
	"github.com/semirm-dev/mahala/internal/redis"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) RedisStorage {
	return RedisStorage{
		redisClient: redisClient,
	}
}

func (r RedisStorage) Store(candidate string, votes []Vote) error {
	return r.redisClient.Add(redis.Item{
		Key:   candidate,
		Value: votes,
	})
}

func (r RedisStorage) Get(candidate string) ([]Vote, error) {
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
