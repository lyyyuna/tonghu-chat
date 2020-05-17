package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

const (
	chanListKey             = "channel.list"
	historyPrefix           = "history"
	chatPrefix              = "chat"
	chatLastSeqPrefix       = "last_seq"
	chatClientLastSeqPrefix = "client.last_seq"

	maxHistorySize int64 = 1000
)

// RedisClient represents the redis client
type RedisClient struct {
	cl *redis.Client
}

// NewRedisClient starts a new redis client
func NewRedisClient(host, pass string, port int) (*RedisClient, error) {
	opts := redis.Options{
		Addr: host + ":" + strconv.Itoa(port),
	}
	if pass != "" {
		opts.Password = pass
	}

	client := redis.NewClient(&opts)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Redis Addr %v, Port %v Reason %v", host, port, err)
	}

	return &RedisClient{cl: client}, nil
}

// Get retrieve a chat from the redis
func (r *RedisClient) Get(id string) {
	val, err := r.cl.Get(chatID(id)).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func chatID(id string) string {
	return fmt.Sprintf("%s.%s", chatPrefix, id)
}

func chatHistoryID(id string) string {
	return fmt.Sprintf("%s.%s.%s", historyPrefix, chatPrefix, id)
}

func chatLastSeqID(id string) string {
	return fmt.Sprintf("%s.%s.%s", chatLastSeqPrefix, chatPrefix, id)
}

func chatClientLastSeqID(uid, id string) string {
	return fmt.Sprintf("%s.%s.%s", chatClientLastSeqPrefix, uid, id)
}
