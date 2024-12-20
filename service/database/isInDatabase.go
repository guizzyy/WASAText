package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) IsInDatabase(uID uint64) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT 1 FROM users WHERE id = ? LIMIT 1`, uID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrUserNotFound
		}
		return false, err
	}
	return true, nil
}
