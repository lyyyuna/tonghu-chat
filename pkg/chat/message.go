package chat

// Message represents chat message
type Message struct {
	Meta     map[string]string `json:"meta"`
	Time     int64             `json:"time"`
	Seq      uint64            `json:"seq"`
	Text     string            `json:"text"`
	FromUID  string            `json:"from_uid"`
	FromName string            `json:"from_name"`
}
