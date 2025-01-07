package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) LogUser(u *utilities.User) (bool, error) {
	err := db.c.QueryRow(`SELECT id, photo FROM user WHERE name = ?`, u.Username).Scan(&u.ID, &u.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := db.c.QueryRow(`INSERT INTO user(name) VALUES (?) RETURNING id, photo`, u.Username).Scan(&u.ID, &u.Photo)
			if err != nil {
				return false, fmt.Errorf("failed to insert a new user: %w", err)
			}
		}
		return false, fmt.Errorf("failed to query user table for login: %w", err)
	}
	return false, nil
}

func (db *appdbimpl) SetUsername(u utilities.User) error {
	res, err := db.c.Exec(`UPDATE user SET name = ? WHERE id = ?`, u.Username, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows in set username db function: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (db *appdbimpl) SetPhoto(u utilities.User) error {
	res, err := db.c.Exec(`UPDATE users SET photo = ? WHERE id = ?`, u.Photo, u.ID)
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

func (db *appdbimpl) GetUsers(username string, id uint64) ([]utilities.User, error) {
	// Get the users wanted with a given username string (avoiding the user self)
	rows, err := db.c.Query(`SELECT name, photo FROM user WHERE id != ? AND name LIKE ?`, id, username+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get users from database: %w", err)
	}
	defer rows.Close()

	// Scan the rows in order to get info about the users found
	var users []utilities.User
	for rows.Next() {
		var user utilities.User
		if err := rows.Scan(&user.Username, &user.Photo); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return users, nil
}

func (db *appdbimpl) GetIDByUsername(username string) (uint64, error) {
	var uID uint64
	err := db.c.QueryRow("SELECT id FROM user WHERE name = ?", username).Scan(&uID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("failed to get the username from database: %w", err)
	}
	return uID, nil
}
