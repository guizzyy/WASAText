package utilities

import "time"

/*
type messStatus string

const (
	StatusReceived   messStatus = "Received"
	StatusRead       messStatus = "Read"
	StatusUnreceived messStatus = "Unreceived"
)
*/

type convType string

const (
	Private convType = "private"
	Group   convType = "group"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

type Message struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Conv      uint64    `json:"conv_id"`
	Sender    uint64    `json:"sender_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type Reaction struct {
	Emoji string `json:"emoji"`
	User  string `json:"user"`
}

type Notification struct {
	Report string `json:"report"`
}

type Conversation struct {
	ID        uint64    `json:"id"`
	Type      convType  `json:"type"`
	Name      string    `json:"name"`
	Photo     string    `json:"photo"`
	LastMess  string    `json:"last_mess"`
	Timestamp time.Time `json:"timestamp"`
}

type LoginResponse struct {
	Message    string `json:"message"`
	UserLogged User   `json:"user"`
}
