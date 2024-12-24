package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) LogUser(u *utilities.User) (bool, error) {
	err := db.c.QueryRow("SELECT id, photo FROM users WHERE name = ?", u.Username).Scan(&u.ID, &u.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res, err := db.c.Exec("INSERT INTO users(name) VALUES (?)", u.Username)
			if err != nil {
				return false, fmt.Errorf("failed to insert u: %w", err)
			}
			id, err := res.LastInsertId()
			if err != nil {
				return false, fmt.Errorf("failed to get the last ID: %w", err)
			}
			u.ID = uint64(id)
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (db *appdbimpl) SetUsername(newName string, id uint64) error {
	res, err := db.c.Exec(`UPDATE users SET name = ? WHERE id = ?`, newName, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (db *appdbimpl) SetPhoto(photo string, id uint64) error {
	res, err := db.c.Exec(`UPDATE users SET photo = ? WHERE id = ?`, photo, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
