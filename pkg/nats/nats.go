package nats

import (
	"encoding/json"
	"fmt"
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
	stan "github.com/nats-io/stan.go"
	"io"
	"time"
)

// NatsClient represents NATS client
type NatsClient struct {
	cn stan.Conn
}

// NewNatsClient initializes a connection to NATS server
func NewNatsClient(clusterID, clientID, url string) (*NatsClient, error) {
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %v", err)
	}
	return &NatsClient{cn: conn}, nil
}

func (nc *NatsClient) Send(id string, msg []byte) error {
	return nc.cn.Publish(id, msg)
}

func (nc *NatsClient) subscribeSeq(id string, start uint64, f func(uint64, []byte)) (stan.Subscription, error) {
	return nc.cn.Subscribe(
		id,
		func(m *stan.Msg) {
			f(m.Sequence, m.Data)
		},
		stan.StartAtSequence(start),
		stan.SetManualAckMode())
}

func (nc *NatsClient) Subscribe(chatId string, uid string, start uint64, c chan *chat.Message) (io.Closer, error) {
	closer, err := nc.subscribeSeq("chat."+chatId, start, func(seq uint64, data []byte) {
		var msg chat.Message
		err := json.Unmarshal(data, &msg)
		if err != nil {
			msg = chat.Message{
				FromUID: "broker",
				Text:    "broker: message unavailable: decoding error",
				Time:    time.Now().UnixNano(),
			}
		}
		msg.Seq = seq

		if msg.FromUID != uid {
			c <- &msg
		} else {

		}
	})

	return closer, err
}
