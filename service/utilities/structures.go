package utilities

type messStatus string

const (
	StatusReceived messStatus = "Received"
	StatusRead     messStatus = "Read"
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
	Text   string `json:"text"`
	Sender uint64 `json:"sender"`
}

type Notification struct {
	Report string `json:"report"`
}

type Conversation struct {
	ConvID uint64 `json:"id"`
	Name   string `json:"name"`
	Photo  string `json:"photo"`
}

type LoginResponse struct {
	Message    string `json:"message"`
	UserLogged User   `json:"user"`
}
