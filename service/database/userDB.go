package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) LogUser(u *utilities.User) (bool, error) {
	//	Select info about the user from the database to figure if it's a new/existing user
	var photo sql.NullString
	err := db.c.QueryRow(`SELECT id, photo FROM user WHERE name = ?`, u.Username).Scan(&u.ID, &photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the user is new, insert him in the database
			err = db.c.QueryRow(`INSERT INTO user(name) VALUES (?) RETURNING id`, u.Username).Scan(&u.ID)
			if err != nil {
				return false, fmt.Errorf("failed to insert a new user: %w", err)
			}
			return true, nil
		}
		return false, fmt.Errorf("failed to query user table for login: %w", err)
	}
	u.Photo = photo.String
	return false, nil
}

func (db *appdbimpl) SetUsername(u utilities.User) error {
	//	Check if the username selected is available
	if isIn, err := db.IsUsernameInDatabase(u.Username); err != nil {
		return fmt.Errorf("failed to check if username is in database: %w", err)
	} else if isIn {
		return errors.New("username is already taken")
	}

	//	Update the database with the new username given and check errors
	res, err := db.c.Exec(`UPDATE user SET name = ? WHERE id = ?`, u.Username, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user name: %w", err)
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
	//	Update the database with the new photo given and check errors
	res, err := db.c.Exec(`UPDATE user SET photo = ? WHERE id = ?`, u.Photo, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user photo: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows in set user photo db function: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (db *appdbimpl) GetUsers(username string, id uint64) ([]utilities.User, error) {
	//	Get the users wanted with a given username string (avoiding the user self)
	rows, err := db.c.Query(`SELECT id, name, photo FROM user WHERE id != ? AND name LIKE ?`, id, username+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get users from database: %w", err)
	}

	//	Scan the rows in order to get info about the users found
	var users []utilities.User
	for rows.Next() {
		var user utilities.User
		var photo sql.NullString
		if err = rows.Scan(&user.ID, &user.Username, &photo); err != nil {
			return nil, fmt.Errorf("error in scanning users info for search: %w", err)
		}
		user.Photo = photo.String
		users = append(users, user)
	}

	//	Check error during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows in GetUsers: %w", err)
	}
	return users, nil
}

func (db *appdbimpl) GetUserByUsername(u *utilities.User) error {
	//	Check if the username exists in the database
	if exists, err := db.IsUsernameInDatabase(u.Username); err != nil {
		return fmt.Errorf("failed to check if username is in database: %w", err)
	} else if !exists {
		return ErrUserNotFound
	}

	//	Get user information with a given username
	var photo sql.NullString
	err := db.c.QueryRow(`SELECT * FROM user WHERE name = ?`, u.Username).Scan(&u.ID, &u.Username, &photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to retrieve user info from db: %w", err)
	}
	u.Photo = photo.String
	return nil
}

func (db *appdbimpl) GetUserByID(uID uint64) (utilities.User, error) {
	//	Get user information with a given user id
	var user utilities.User
	var photo sql.NullString
	err := db.c.QueryRow(`SELECT * FROM user WHERE id = ?`, uID).Scan(&user.ID, &user.Username, &photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("failed to retrieve user from db: %w", err)
	}
	user.Photo = photo.String
	return user, nil
}
