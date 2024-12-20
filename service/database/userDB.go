package database

func (db *appdbimpl) LogUser(username string) (DUser, error) {
	var u DUser
	res, err := db.c.Exec(`INSERT INTO users(username) VALUES (?)`, username)
	if err != nil {
		check, _ := res.RowsAffected()
		if check == 0 {
			if err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&u.Id); err != nil {
				return u, err
			}
			return u, nil
		}
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = uint64(lastInsertID)
	u.Username = username
	return u, nil
}

func (db *appdbimpl) SetUsername(newName string, id uint64) error {
	res, err := db.c.Exec(`UPDATE users SET username = ? WHERE id = ?`, newName, id)
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
