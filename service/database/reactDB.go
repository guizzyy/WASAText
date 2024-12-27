package database

import (
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) AddReaction(react utilities.Reaction, messId uint64) error {
	_, err := db.c.Exec(`INSERT INTO reactions(reaction, mess_id, sender) VALUES (?, ?, ?)`, react.Emoji, messId, react.User)
	if err != nil {
		return fmt.Errorf("error adding reaction: %s", err)
	}
	return nil
}

func (db *appdbimpl) RemoveReaction(messId uint64, sender string) error {
	_, err := db.c.Exec(`DELETE FROM reactions WHERE mess_id = ? AND sender = ?`, messId, sender)
	if err != nil {
		return fmt.Errorf("error removing reaction: %s", err)
	}
	return nil
}
