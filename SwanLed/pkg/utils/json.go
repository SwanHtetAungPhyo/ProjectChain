package utils

import "encoding/json"

// Message represents a JSON message structure
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// Encode encodes a Message into JSON format
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode decodes JSON data into a Message struct
func Decode(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
