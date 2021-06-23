package relay

import (
	"encoding/json"
)

type MessageType string

const (
	MessageTypeStart = "start"
	MessageTypeEvent = "event"
)

type IncomingMessageStart struct {
	Type MessageType              `json:"type"`
	Data IncomingMessageStartData `json:"data"`
}

type IncomingMessageStartData struct {
	ConnectionID string `json:"connection_id"`
}

type IncomingMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type IncomingMessageEventData struct {
	ID      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type OutgoingMessageEvent struct {
	Type MessageType              `json:"type"`
	Data OutgoingMessageEventData `json:"data"`
}

type OutgoingMessageEventData struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
