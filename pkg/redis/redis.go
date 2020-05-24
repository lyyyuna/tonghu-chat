package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
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

// GetChannel retrieve a channel from the redis
func (r *RedisClient) GetChannel(id string) (*chat.Channel, error) {
	val, err := r.cl.Get(chatID(id)).Result()
	if err != nil {
		return nil, err
	}

	var ch chat.Channel
	err = json.Unmarshal([]byte(val), &ch)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (r *RedisClient) SaveChannel(ch *chat.Channel) error {
	data, err := json.Marshal(ch)
	if err != nil {
		return err
	}

	pipe := r.cl.TxPipeline()
	pipe.Set(chatID(chatID(ch.Name)), data, 0)

	_, err = pipe.Exec()
	return err
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
