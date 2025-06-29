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
	GetUserByUsername(user *utilities.User) error
	GetUserByID(uint64) (utilities.User, error)

	SetGroupName(utilities.Conversation, uint64) error
	SetGroupPhoto(utilities.Conversation, uint64) error
	CreatePrivConv(utilities.User, utilities.User) (utilities.Conversation, error)
	CreateGroupConv(*utilities.Conversation, uint64) error
	AddToGroup(uint64, uint64, utilities.User) error
	LeaveGroup(uint64, uint64) error
	IsGroupConv(uint64) (bool, error)
	IsUserInConv(uint64, uint64) (bool, error)
	PrivConvExists(utilities.User, utilities.User) (bool, utilities.Conversation, error)
	GroupStillExists(uint64) error

	GetConversations(uint64) ([]utilities.Conversation, error)
	GetConversation(utilities.Conversation, uint64) ([]utilities.Message, error)
	GetReceivers(uint64, uint64) ([]uint64, error)
	GetMembers(uint64, uint64) ([]utilities.User, error)
	GetGroupPhoto(uint64) (string, error)
	GetPrivConvInfo(uint64, uint64) (string, string, error)
	GetGroupConvInfo(uint64) (string, string, error)
	GetConvByID(uint64, uint64) (utilities.Conversation, error)

	GetMessageInfo(uint64) (utilities.Message, error)
	AddMessage(*utilities.Message) error
	RemoveMessage(uint64, uint64) error
	GetLastMessage(utilities.Conversation, uint64) (utilities.Message, error)
	InsertStatus([]uint64, uint64, uint64) (string, error)
	UpdateReceivedStatus(uint64) error
	UpdateReadStatus(uint64, uint64) error
	CheckStatus(uint64, uint64) (string, error)
	IsOwnerMessage(uint64, uint64) (bool, error)
	IsMessageInConv(uint64, uint64) (bool, error)

	AddReaction(utilities.Reaction, uint64) error
	RemoveReaction(uint64, uint64) error
	GetReactions(uint64) ([]utilities.Reaction, error)

	Ping() error
	IsUserInDatabase(uint64) (bool, error)
	IsConvInDatabase(uint64) (bool, error)
	IsMessageInDatabase(uint64) (bool, error)
	IsReactionInDatabase(uint64, uint64) (bool, error)
	IsMembershipInDatabase(uint64, uint64) (bool, error)
	IsUsernameInDatabase(string) (bool, error)
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
			"user": `CREATE TABLE IF NOT EXISTS user (
    		id INTEGER PRIMARY KEY, 
    		name VARCHAR(16) COLLATE NOCASE UNIQUE NOT NULL CHECK ( length(name) >= 3 AND length(name) <= 16 ),
    		photo TEXT DEFAULT NULL)`,

			"conversation": `CREATE TABLE IF NOT EXISTS conversation (
    		id INTEGER PRIMARY KEY,
    		type TEXT CHECK ( type IN ('private', 'group') ),
    		name VARCHAR(25) CHECK ( length(name) >= 3 AND length(name) <= 25 ),
    		photo TEXT DEFAULT NULL)`,

			"membership": `CREATE TABLE IF NOT EXISTS membership (
    		conv_id INTEGER NOT NULL,
    		user_id INTEGER NOT NULL,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		UNIQUE (conv_id, user_id),
    		FOREIGN KEY (conv_id) REFERENCES conversation(id) ON DELETE CASCADE,
    		FOREIGN KEY (user_id) REFERENCES user(id))`,

			"message": `CREATE TABLE IF NOT EXISTS message (
    		id INTEGER PRIMARY KEY,
    		text VARCHAR(250) CHECK ( length(text) >= 0 AND length(text) <= 250 ),
    		photo TEXT DEFAULT NULL,
    		conv_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		reply_to INTEGER NOT NULL DEFAULT 0,
    		is_forwarded BOOLEAN NOT NULL DEFAULT FALSE,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		FOREIGN KEY (conv_id) REFERENCES conversation(id) ON DELETE CASCADE,
    		FOREIGN KEY (sender_id) REFERENCES user(id))`,

			"status": `CREATE TABLE IF NOT EXISTS status (
    		receiver_id INTEGER NOT NULL,
    		mess_id INTEGER NOT NULL,
    		conv_id INTEGER NOT NULL,
    		info TEXT DEFAULT 'Unreceived' CHECK ( info IN ('Read', 'Received', 'Unreceived') ),
    		FOREIGN KEY (mess_id) REFERENCES message(id) ON DELETE CASCADE,
    		FOREIGN KEY (receiver_id) REFERENCES user(id),
    		FOREIGN KEY (conv_id) REFERENCES conversation(id) ON DELETE CASCADE,
    		PRIMARY KEY (mess_id, receiver_id, conv_id))`,

			"reactions": `CREATE TABLE IF NOT EXISTS reactions (
    		reaction TEXT NOT NULL,
    		mess_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		PRIMARY KEY (mess_id, sender_id, reaction),
    		FOREIGN KEY (mess_id) REFERENCES message(id) ON DELETE CASCADE,
    		FOREIGN KEY (sender_id) REFERENCES user(id))`,
		}

		for table, query := range tables {
			if _, err = db.Exec(query); err != nil {
				return nil, fmt.Errorf("error creating table %s: %w", table, err)
			}
		}

	}

	return &appdbimpl{
		c: db,
	}, nil
}

var ErrUserNotFound = errors.New("user not found")
var ErrMessageNotFound = errors.New("message not found")
var ErrConversationNotFound = errors.New("conversation not found")
var ErrMembershipNotFound = errors.New("membership not found")
var ErrReactionNotFound = errors.New("reaction not found")
var ErrNoSelfConversation = errors.New("can't create self conversation")
var ErrNoGroup = errors.New("conversation is not a group")
var ErrUserInGroup = errors.New("user is already in this group")
var ErrUserNotInConversation = errors.New("user is not in this conversation")

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
