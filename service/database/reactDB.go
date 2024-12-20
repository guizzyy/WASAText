package database

import "fmt"

func (db *appdbimpl) AddReaction(emoji string, messId uint64, senderId uint64) error {
	_, err := db.c.Exec(`INSERT INTO reacts_messages(reaction, messId, sender_id) VALUES ($1, $2, $3)`, emoji, messId, senderId)
	if err != nil {
		return fmt.Errorf("error adding reaction: %s", err)
	}
	return nil
}

func (db *appdbimpl) RemoveReaction(messId uint64, senderId uint64) error {
	_, err := db.c.Exec(`DELETE FROM reacts_messages WHERE mess_id = $1 AND sender_id = $2`, messId, senderId)
	if err != nil {
		return fmt.Errorf("error removing reaction: %s", err)
	}
	return nil
}
