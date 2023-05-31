package voting

import "github.com/sirupsen/logrus"

type RedisStorage struct{}

func (r RedisStorage) Store(vote Vote) error {
	logrus.Infof("vote %v stored in redis", vote)
	return nil
}
