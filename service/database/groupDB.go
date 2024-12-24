package database

import (
	"fmt"
)

func (db *appdbimpl) SetGroupName(groupName string, groupId uint64) error {
	res, err := db.c.Exec(`UPDATE groups SET name = ? WHERE id = ?`, groupName, groupId)
	if err != nil {
		return fmt.Errorf("error updating group name: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (db *appdbimpl) SetGroupPhoto(groupId uint64, photo string) error {
	res, err := db.c.Exec(`UPDATE groups SET photo = ? WHERE id = ?`, photo, groupId)
	if err != nil {
		return fmt.Errorf("error updating group photo: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (db *appdbimpl) AddMembership(userId uint64, groupId uint64) error {
	_, err := db.c.Exec(`INSERT INTO membership(group_id, user_id) VALUES (?, ?)`, groupId, userId)
	if err != nil {
		return fmt.Errorf("error adding membership: %v", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMembership(userId uint64, groupId uint64) error {
	_, err := db.c.Exec(`DELETE FROM membership WHERE group_id = ? AND user_id = ?`, groupId, userId)
	if err != nil {
		return fmt.Errorf("error removing membership: %v", err)
	}
	return nil
}
