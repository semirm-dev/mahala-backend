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

func (r RedisStorage) Store(vote Vote) error {
	candidateVotes, err := r.redisClient.Get(vote.Candidate)
	if err != nil && err != redis.ErrNotExists {
		return err
	}

	var votes []Vote
	if len(candidateVotes) > 0 {
		if err := json.Unmarshal(candidateVotes, &votes); err != nil {
			return err
		}
	}

	votes = append(votes, Vote{
		Candidate: vote.Candidate,
		VoterID:   vote.VoterID,
	})

	return r.redisClient.Add(redis.Item{
		Key:   vote.Candidate,
		Value: votes,
	})
}
