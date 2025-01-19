package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) IsUserInDatabase(uID uint64) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM user WHERE id = ? LIMIT 1`, uID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrUserNotFound
		}
		return false, fmt.Errorf("error during IsUserInDatabase: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IsConvInDatabase(cID uint64) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM conversation WHERE id = ? LIMIT 1`, cID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrConversationNotFound
		}
		return false, fmt.Errorf("error during IsConvInDatabase: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IsMessageInDatabase(mID uint64) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM message WHERE id = ? LIMIT 1`, mID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrMessageNotFound
		}
		return false, fmt.Errorf("error during IsMessageInDatabase: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IsReactionInDatabase(react string, mID uint64, uID uint64) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM reactions WHERE (reaction, mess_id, sender_id) = (?, ?, ?) LIMIT 1`, react, mID, uID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrReactionNotFound
		}
		return false, fmt.Errorf("error during IsReactionInDatabase: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IsMembershipInDatabase(uID uint64, cID uint64) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM membership WHERE (conv_id, user_id) = (?, ?) LIMIT 1`, cID, uID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrMembershipNotFound
		}
		return false, fmt.Errorf("error during IsMembershipInDatabase: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IsUsernameInDatabase(username string) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM user WHERE name = ? LIMIT 1`, username).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error during IsUsernameInDatabase: %w", err)
	}
	return count > 0, nil
}
