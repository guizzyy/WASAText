package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) GetConversations(uID uint64) ([]utilities.Conversation, error) {
	// TODO: change the query in order to get last message and ordered it (also manage the status)

	// Select the conv infos where the user participates
	rows, err := db.c.Query(`SELECT (conversation.id, type, name, photo) FROM message INNER JOIN conversation INNER JOIN memberships ON conversation.id = memberships.conv_id = message.conv_id WHERE user_id = ?`, uID)
	if err != nil {
		return nil, fmt.Errorf("error in getting conversations info: %w", err)
	}
	defer rows.Close()

	// Create an array of conversation structs to return and scan the rows
	convs := make([]utilities.Conversation, 0)
	for rows.Next() {
		var conv utilities.Conversation
		if err := rows.Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo); err != nil {
			return nil, fmt.Errorf("error in scanning conversation info: %w", err)
		}
		convs = append(convs, conv)
	}

	// Check errors during the scan, otherwise return the array
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}
	return convs, nil
}

func (db *appdbimpl) GetConversation(convID uint64) ([]utilities.Message, error) {
	// TODO: find a way to manage the status of the message situation
	rows, err := db.c.Query(`SELECT id, text, sender_id, timestamp FROM message WHERE conv_id = ?`, convID)
}

func (db *appdbimpl) CreateGroupConv(grName string) (utilities.Conversation, error) {
	var conv utilities.Conversation
	res := db.c.QueryRow(`INSERT INTO conversation(name) VALUES (?) RETURNING *`, grName).Scan()
	// TODO: check query correctness and continue
}

func (db *appdbimpl) SetGroupName(group *utilities.Conversation) error {
	_ := db.c.QueryRow(`UPDATE conversation SET name = ? WHERE id = ? RETURNING *`, group.Name, group.ID).Scan()
	// TODO: same as the function above
}

func (db *appdbimpl) SetGroupPhoto(group *utilities.Conversation) error {
	_ := db.c.QueryRow(`UPDATE conversation SET photo ? WHERE id = ? RETURNING *`, group.Photo, group.ID).Scan()
	// TODO: same as the function above
}

func (db *appdbimpl) AddToGroup(idConv uint64, u utilities.User) error {
	_, err := db.c.Exec(`INSERT INTO memberships(conv_id, user_id) VALUES (?, ?)`, idConv, u.ID)
	if err != nil {
		return fmt.Errorf("error in adding membership to conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(idConv uint64, idUser uint64) error {
	_, err := db.c.Exec(`DELETE FROM memberships WHERE conv_id = ? AND user_id = ?`, idConv, idUser)
	if err != nil {
		return fmt.Errorf("error in leaving conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetReceiver(convID uint64, senderID uint64) (uint64, error) {
	var receiver uint64
	err := db.c.QueryRow(`SELECT user_id FROM memberships WHERE conv_id = ? AND user_id != ?`, convID, senderID).Scan(&receiver)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrMembershipNotFound
		}
	}
	return receiver, nil
}
