package database

import (
	"fmt"
)

func (db *appdbimpl) AddMessage(text string, convId uint64, senderId uint64) error {
	res, err := db.c.Exec(`INSERT INTO messages (text, conv_id, sender_id) VALUES (?, ?, ?)`, text, convId, senderId)
	if err != nil {
		return fmt.Errorf("error adding message to database: %v", err)
	}

	idMess, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}
	if _, err := db.c.Exec(`INSERT INTO status_messages (idMess, receiver_id) VALUES (?, ?)`, idMess, receiverId); err != nil {
		return fmt.Errorf("error adding the message status to database: %v", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMessage(messId uint64) error {
	_, err := db.c.Exec(`DELETE FROM messages WHERE idMess = ?`, messId)
	if err != nil {
		return fmt.Errorf("error removing the message from the database: %v", err)
	}
	return nil
}

func (db *appdbimpl) ForwardMessage(messId uint64, receiverId uint64) error {
	var text string
	var convId uint64
	var senderId uint64
	rows := db.c.QueryRow(`SELECT (text, conv_id, sender_id) FROM messages WHERE id = ?`, messId).Scan(&text, &convId, &senderId)
	if rows == nil {
		return fmt.Errorf("could not find message to forward to")
	}
	// TO DO: CONTINUE
}
