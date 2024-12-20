package api

type messStatus string

const (
	StatusReceived messStatus = "Received"
	StatusRead     messStatus = "Read"
)

type Username struct {
	Username string `json:"username"`
}

type Photo struct {
	Photo string `json:"photo"`
}

type Notification struct {
	Outcome   bool   `json:"outcome"`
	Report    string `json:"report"`
	ErrorCode int    `json:"error"`
}
