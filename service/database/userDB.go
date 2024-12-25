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

func (db *appdbimpl) SetUsername(u utilities.User) error {
	res, err := db.c.Exec(`UPDATE users SET name = ? WHERE id = ?`, u.Username, u.ID)
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
