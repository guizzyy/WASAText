/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	LogUser(*utilities.User) (bool, error)
	SetUsername(utilities.User) error
	SetPhoto(utilities.User) error
	GetUsers(string, uint64) ([]utilities.User, error)
	GetIDByUsername(string) (uint64, error)

	SetGroupName(utilities.Conversation) error
	SetGroupPhoto(utilities.Conversation) error
	CreateGroupConv(*utilities.Conversation, uint64) error
	AddToGroup(uint64, utilities.User) error
	LeaveGroup(uint64, uint64) error
	IsGroupConv(uint64) (bool, error)

	GetConversations(uint64) ([]utilities.Conversation, error)
	GetConversation(uint64, uint64) ([]utilities.Message, error)
	GetReceivers(uint64, uint64) ([]uint64, error)

	GetMessageInfo(uint64) (utilities.Message, error)
	AddMessage(*utilities.Message) error
	RemoveMessage(uint64) error
	InsertStatus([]uint64, uint64) (string, error)
	UpdateReceivedStatus(*utilities.Message) error
	UpdateReadStatus(*utilities.Message) error

	AddReaction(utilities.Reaction, uint64) error
	RemoveReaction(uint64, uint64) error
	GetReactions(uint64) ([]utilities.Reaction, error)

	Ping() error
	IsUserInDatabase(uint64) (bool, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable foreign keys in the database
	_, errFK := db.Exec("PRAGMA foreign_keys = ON;")
	if errFK != nil {
		return &appdbimpl{}, errFK
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		tables := map[string]string{
			"user": `CREATE TABLE user (
    		id INTEGER PRIMARY KEY, 
    		name VARCHAR(16) UNIQUE NOT NULL CHECK ( length(name) >= 3 AND length(name) <= 16 ),
    		photo TEXT DEFAULT NULL)`,

			"conversation": `CREATE TABLE conversation (
    		id INTEGER PRIMARY KEY,
    		type TEXT CHECK ( type IN ('private', 'group') ),
    		name VARCHAR(25) NOT NULL CHECK ( length(name) >= 3 AND length(name) <= 25 ))`,

			"membership": `CREATE TABLE membership (
    		conv_id INTEGER NOT NULL,
    		user_id INTEGER NOT NULL,
    		UNIQUE (conv_id, user_id),
    		FOREIGN KEY (conv_id) REFERENCES conversation(id),
    		FOREIGN KEY (user_id) REFERENCES user(id))`,

			"message": `CREATE TABLE message (
    		id INTEGER PRIMARY KEY,
    		text VARCHAR(250) CHECK ( length(text) >= 0 AND length(text) <= 250 ),
    		photo TEXT DEFAULT NULL,
    		conv_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		is_forwarded BOOLEAN NOT NULL DEFAULT FALSE,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		FOREIGN KEY (conv_id) REFERENCES conversations(id),
    		FOREIGN KEY (sender_id) REFERENCES users(id))`,

			"status": `CREATE TABLE status (
    		receiver_id INTEGER NOT NULL,
    		mess_id INTEGER NOT NULL,
    		info TEXT DEFAULT 'Unreceived' CHECK ( info IN ('Read', 'Received', 'Unreceived') ),
    		FOREIGN KEY (mess_id) REFERENCES messages(id) ON DELETE CASCADE,
    		FOREIGN KEY (receiver_id) REFERENCES users(id),
    		PRIMARY KEY (mess_id, receiver_id))`,

			"reactions": `CREATE TABLE reactions (
    		reaction TEXT NOT NULL,
    		mess_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		PRIMARY KEY (mess_id, sender_id, reaction),
    		FOREIGN KEY (mess_id) REFERENCES messages(id) ON DELETE CASCADE,
    		FOREIGN KEY (sender_id) REFERENCES users(id))`,
		}

		for table, query := range tables {
			if _, err := db.Exec(query); err != nil {
				return nil, fmt.Errorf("error creating table %s: %v", table, err)
			}
		}

	}

	return &appdbimpl{
		c: db,
	}, nil
}

var ErrUserNotFound = errors.New("User not found")
var ErrMessageNotFound = errors.New("Message not found")
var ErrConversationNotFound = errors.New("Conversation not found")
var ErrMembershipNotFound = errors.New("Membership not found")
var ErrReactionNotFound = errors.New("Reaction not found")

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
