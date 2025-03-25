package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) AddReaction(react utilities.Reaction, messId uint64) error {
	//	Check if the message exists in order to comment it
	if exists, err := db.IsMessageInDatabase(messId); err != nil {
		return fmt.Errorf("error in checking if message exists: %w", err)
	} else if !exists {
		return ErrMessageNotFound
	}

	//	Check if already exists a reaction of the user in a message
	if isIn, err := db.IsReactionInDatabase(messId, react.User.ID); err != nil {
		return fmt.Errorf("error in checking if the reaction exists: %w", err)
	} else if !isIn {
		//	if it doesn't exist, insert a new one
		_, err = db.c.Exec(`INSERT INTO reactions(reaction, mess_id, sender_id) VALUES (?, ?, ?)`, react.Emoji, messId, react.User.ID)
		if err != nil {
			return fmt.Errorf("error adding reaction: %w", err)
		}
		return nil
	} else {
		//	if it exists, update it with another reaction
		_, err = db.c.Exec(`UPDATE reactions SET reaction = ?, timestamp = CURRENT_TIMESTAMP WHERE sender_id = ? AND mess_id = ?`, react.Emoji, react.User.ID, messId)
		if err != nil {
			return fmt.Errorf("error updating reaction: %w", err)
		}
		return nil
	}
}

func (db *appdbimpl) RemoveReaction(messId uint64, senderId uint64) error {
	//	Check if the reaction of the user is in the database
	if isIn, err := db.IsReactionInDatabase(messId, senderId); err != nil {
		return fmt.Errorf("error in checking if reaction exists: %w", err)
	} else if !isIn {
		return ErrReactionNotFound
	}

	//	Remove the reaction from a message (only if you are the sender)
	_, err := db.c.Exec(`DELETE FROM reactions WHERE mess_id = ? AND sender_id = ?`, messId, senderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReactionNotFound
		}
		return fmt.Errorf("error removing reaction: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetReactions(messId uint64) ([]utilities.Reaction, error) {
	//	Get the reaction emojis and senders from the database
	var reactions []utilities.Reaction
	query := `SELECT 
					r.reaction,
					u.name,
					u.photo
				FROM 
				    reactions AS r, 
				    user AS u
				WHERE
				    r.mess_id = ? AND
				    r.sender_id = u.id
    `
	rows, err := db.c.Query(query, messId)
	if err != nil {
		return nil, fmt.Errorf("error in getting reactions of the message: %w", err)
	}

	//	Iterate the rows to save information in the array
	for rows.Next() {
		var reaction utilities.Reaction
		var sender utilities.User
		if err = rows.Scan(&reaction.Emoji, &sender.Username, &sender.Photo); err != nil {
			return nil, fmt.Errorf("error in scanning reactions of the message: %w", err)
		}
		reaction.User = sender
		reactions = append(reactions, reaction)
	}

	//	Check error during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error resulting rows of GetReactions: %w", err)
	}

	return reactions, nil
}
