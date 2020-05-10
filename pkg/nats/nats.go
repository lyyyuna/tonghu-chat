package nats

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"io"
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

func (nc *NatsClient) SubscribeSeq(id string, start uint64, f func(uint64, []byte)) (stan.Subscription, error) {
	return nc.cn.Subscribe(
		id,
		func(m *stan.Msg) {
			f(m.Sequence, m.Data)
		},
		stan.StartAtSequence(start),
		stan.SetManualAckMode())
}
