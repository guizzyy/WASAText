package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) SetGroupName(groupName string, groupId uint64) error {
	res, err := db.c.Exec(`UPDATE groups SET groupName = ? WHERE groupId = ?`, groupName, groupId)
	if err != nil {
		return fmt.Errorf("error updating group name: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("group does not exist")
	}
	return nil
}

func (db *appdbimpl) SetGroupPhoto(groupId uint64, photo string) error {
	res, err := db.c.Exec(`UPDATE groups SET photo = ? WHERE groupId = ?`, photo, groupId)
	if err != nil {
		return fmt.Errorf("error updating group photo: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("group does not exist")
	}
	return nil
}

func (db *appdbimpl) AddMembership(userId uint64, groupId uint64) error {
	_, err := db.c.Exec(`INSERT INTO memberships (groupId, userId) VALUES (?, ?)`, groupId, userId)
	if err != nil {
		return fmt.Errorf("error adding membership: %v", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMembership(userId uint64, groupId uint64) error {
	_, err := db.c.Exec(`DELETE FROM memberships WHERE groupId = ? AND userId = ?`, groupId, userId)
	if err != nil {
		return fmt.Errorf("error removing membership: %v", err)
	}
	return nil
}
