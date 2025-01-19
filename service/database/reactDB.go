package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) AddReaction(react utilities.Reaction, messId uint64) error {
	if isIn, err := db.IsReactionInDatabase(react.Emoji, messId, react.User); err != nil {
		return fmt.Errorf("error in checking if the reaction exists: %w", err)
	} else if !isIn {
		_, err = db.c.Exec(`INSERT INTO reactions(reaction, mess_id, sender_id) VALUES (?, ?, ?)`, react.Emoji, messId, react.User)
		if err != nil {
			return fmt.Errorf("error adding reaction: %s", err)
		}
		return nil
	} else {
		_, err = db.c.Exec(`UPDATE reactions SET reaction = ?, timestamp = CURRENT_TIMESTAMP WHERE sender_id = ? AND mess_id = ?`, react.Emoji, react.User, messId)
		if err != nil {
			return fmt.Errorf("error updating reaction: %s", err)
		}
		return nil
	}
}

func (db *appdbimpl) RemoveReaction(messId uint64, senderId uint64) error {
	_, err := db.c.Exec(`DELETE FROM reactions WHERE mess_id = ? AND sender_id = ?`, messId, senderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReactionNotFound
		}
		return fmt.Errorf("error removing reaction: %s", err)
	}
	return nil
}

func (db *appdbimpl) GetReactions(messId uint64) ([]utilities.Reaction, error) {
	var reactions []utilities.Reaction
	// Get the reaction emojis and senders from the database
	rows, err := db.c.Query(`SELECT reaction, sender_id FROM reactions WHERE mess_id = ?`, messId)
	if err != nil {
		return nil, fmt.Errorf("error in getting reactions of the message: %s", err)
	}
	defer rows.Close()

	// Iterate the rows to save information in the array
	for rows.Next() {
		var reaction utilities.Reaction
		if err = rows.Scan(&reaction.Emoji, &reaction.User); err != nil {
			return nil, fmt.Errorf("error in scanning reactions of the message: %s", err)
		}
		reactions = append(reactions, reaction)
	}

	// Check error during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error resulting rows of GetReactions: %s", err)
	}

	return reactions, nil
}
