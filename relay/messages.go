package relay

import (
	"encoding/json"
)

type MessageType string

const version = 1

const (
	MessageTypeStart = "start"
	MessageTypeEvent = "event"
)

type IncomingMessageStart struct {
	Type    MessageType              `json:"type"`
	Version int                      `json:"version"`
	Data    IncomingMessageStartData `json:"data"`
}

type IncomingMessageStartData struct {
	Token string `json:"token"`
}

type IncomingMessage struct {
	Type    MessageType     `json:"type"`
	Version int             `json:"version"`
	Data    json.RawMessage `json:"data"`
}

type IncomingMessageEventData struct {
	ID      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Method  string            `json:"method"`
}

type OutgoingMessageStart struct {
	Type    MessageType              `json:"type"`
	Version int                      `json:"version"`
	Data    OutgoingMessageStartData `json:"data"`
}

type OutgoingMessageStartData struct {
	Token string `json:"token"`
}

type OutgoingMessageEvent struct {
	Type    MessageType              `json:"type"`
	Version int                      `json:"version"`
	Data    OutgoingMessageEventData `json:"data"`
}

type OutgoingMessageEventData struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
