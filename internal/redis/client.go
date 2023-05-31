package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redisLib "github.com/go-redis/redis"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort = 6379
	defaultDb   = 0
)

var ErrNotExists = errors.New("key does not exist")

type Client struct {
	libClient *redisLib.Client
	conf      Config
}

type Config struct {
	Host       string
	Port       int
	Password   string
	DB         int
	PipeLength int
}

type Item struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}

func NewConfig() Config {
	host := os.Getenv("REDIS_HOST")

	port := defaultPort
	portEnv, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err == nil && portEnv > 0 {
		port = portEnv
	}

	password := os.Getenv("REDIS_PASSWORD")

	db := defaultDb
	dbEnv, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err == nil && dbEnv > 0 {
		db = dbEnv
	}

	return Config{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}

func NewClient(conf Config) *Client {
	libClient := redisLib.NewClient(&redisLib.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	return &Client{
		libClient: libClient,
		conf:      conf,
	}
}

func (client *Client) Connect(ctx context.Context) error {
	_, err := client.libClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) Add(item Item) error {
	itemBytes, err := json.Marshal(item.Value)
	if err != nil {
		return err
	}

	return client.libClient.Set(item.Key, string(itemBytes), item.Expiration).Err()
}

func (client *Client) Get(key string) ([]byte, error) {
	val, err := client.libClient.Get(key).Result()

	switch {
	// key does not exist
	case err == redisLib.Nil:
		return nil, ErrNotExists
	// some other error
	case err != nil:
		return nil, err
	}

	return []byte(val), nil
}
