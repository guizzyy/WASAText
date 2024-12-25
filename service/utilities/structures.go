package utilities

import "time"

type messStatus string

const (
	StatusReceived messStatus = "Received"
	StatusRead     messStatus = "Read"
)

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

type ID struct {
	ID string `json:"id"`
}

type Message struct {
	Text      string     `json:"text"`
	Sender    uint64     `json:"sender"`
	Timestamp time.Time  `json:"timestamp"`
	Status    messStatus `json:"status"`
}

type Notification struct {
	Report string `json:"report"`
}

type Conversation struct {
	ID    uint64   `json:"id"`
	Type  convType `json:"type"`
	Name  string   `json:"name"`
	Photo string   `json:"photo"`
}

type LoginResponse struct {
	Message    string `json:"message"`
	UserLogged User   `json:"user"`
}
