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

func (db *appdbimpl) RemoveReaction(messId uint64, senderId uint64) error {
	_, err := db.c.Exec(`DELETE FROM reactions WHERE mess_id = ? AND sender = ?`, messId, senderId)
	if err != nil {
		return fmt.Errorf("error removing reaction: %s", err)
	}
	return nil
}

func (db *appdbimpl) GetReactions(messId uint64) ([]utilities.Reaction, error) {
	var reactions []utilities.Reaction
	// Get the reaction emojis and senders from the database
	rows, err := db.c.Query(`SELECT reaction, sender_id FROM reactions WHERE mess_id = ?`, messId)
	if err != nil {
		return nil, fmt.Errorf("error getting reactions: %s", err)
	}
	defer rows.Close()

	// Iterate the rows to save information in the array
	for rows.Next() {
		var reaction utilities.Reaction
		if err = rows.Scan(&reaction.Emoji, &reaction.User); err != nil {
			return nil, fmt.Errorf("error in scanning the rows: %s", err)
		}
		reactions = append(reactions, reaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in iterate the rows: %s", err)
	}

	return reactions, nil
}
