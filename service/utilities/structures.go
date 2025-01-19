package utilities

import (
	"errors"
	"time"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

type Message struct {
	ID        uint64    `json:"id"`
	Text      string    `json:"text"`
	Photo     string    `json:"photo"`
	Conv      uint64    `json:"conv_id"`
	Sender    uint64    `json:"sender_id"`
	IsForward bool      `json:"is_forwarded"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type Reaction struct {
	Emoji string `json:"emoji"`
	User  uint64 `json:"user"`
}

type Notification struct {
	Report string `json:"report"`
}

type Conversation struct {
	ID          uint64  `json:"id"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Photo       string  `json:"photo"`
	LastMessage Message `json:"last_message"`
}

type LoginResponse struct {
	Message    string `json:"message"`
	UserLogged User   `json:"user"`
}

var ErrNameString = errors.New("invalid string format for the name (length should be between 3 and 25)")
var ErrTextString = errors.New("invalid string format for the message (length should be between 1 and 250)")
