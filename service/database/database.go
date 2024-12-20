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
)

type DUser struct {
	Username string
	Photo    string
	Id       uint64
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	LogUser(string) (DUser, error)
	SetUsername(string, uint64) error
	SetPhoto(string, uint64) error

	Ping() error
	IsInDatabase(uint64) (bool, error)
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

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		tables := map[string]string{
			"users": `CREATE TABLE users (
    		id INTEGER PRIMARY KEY, 
    		name VARCHAR(16) NOT NULL CHECK ( length(name) >= 3 AND length(name) <= 16 ),
    		photo TEXT DEFAULT NULL)`,

			"groups": `CREATE TABLE groups (
    		id INTEGER PRIMARY KEY,
    		name VARCHAR(16) CHECK ( length(name) >= 3 AND length(name) <= 16 ),
    		photo TEXT DEFAULT NULL)`,

			"membership": `CREATE TABLE membership (
    		group_id INTEGER NOT NULL,
    		user_id INTEGER NOT NULL,
    		timestamp_in TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		timestamp_out TIMESTAMP DEFAULT NULL,               
    		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    		FOREIGN KEY (user_id) REFERENCES users(id),
    		PRIMARY KEY (group_id, user_id))`,

			"conversations": `CREATE TABLE conversations (
    		id INTEGER PRIMARY KEY,
    		type TEXT CHECK ( type IN ('Private', 'Group') ),
    		user1_id INTEGER NOT NULL,
    		user2_id INTEGER,
    		group_id INTEGER,
    		FOREIGN KEY (user1_id) REFERENCES users(id),
    		FOREIGN KEY (user2_id) REFERENCES users(id),
    		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE)`,

			"messages": `CREATE TABLE messages (
    		id INTEGER PRIMARY KEY,
    		text TEXT NOT NULL,
    		conv_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		FOREIGN KEY (conv_id) REFERENCES conversations(id),
    		FOREIGN KEY (sender_id) REFERENCES users(id))`,

			"reacts_messages": `CREATE TABLE reacts_messages (
    		reaction TEXT NOT NULL,
    		mess_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		FOREIGN KEY (mess_id) REFERENCES messages(id) ON DELETE CASCADE,
    		FOREIGN KEY (sender_id) REFERENCES users(id),
    		PRIMARY KEY (mess_id, sender_id))`,

			"status_messages": `CREATE TABLE status_messages (
    		mess_id INTEGER NOT NULL,
    		receiver_id INTEGER NOT NULL,
    		is_received BOOLEAN NOT NULL DEFAULT 0,
    		is_read BOOLEAN NOT NULL DEFAULT 0,
    		timestamp_received TIMESTAMP DEFAULT NULL,
    		timestamp_read TIMESTAMP DEFAULT NULL,
    		FOREIGN KEY (mess_id) REFERENCES messages(id) ON DELETE CASCADE,
    		FOREIGN KEY (receiver_id) REFERENCES users(id),
    		PRIMARY KEY (mess_id, receiver_id))`,
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

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
