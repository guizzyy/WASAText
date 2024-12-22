package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetConversations(uID uint64) ([]uint64, error) {
	return nil, nil
}

func (db *appdbimpl) GetConversation(convID uint64) (uint64, error) {
	return 0, nil
}

func (db *appdbimpl) GetReceiver(convID uint64) (uint64, error) {
	var receiver uint64
	err := db.c.QueryRow(`SELECT user2_id FROM conversations WHERE id = ?`, convID).Scan(&receiver)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("conversation not found")
	}
	return receiver, err
}
