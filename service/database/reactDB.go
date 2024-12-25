package database

import "fmt"

func (db *appdbimpl) AddReaction(emoji string, messId uint64, sender string) error {
	_, err := db.c.Exec(`INSERT INTO reactions(reaction, mess_id, sender) VALUES (?, ?, ?)`, emoji, messId, sender)
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
