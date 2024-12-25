package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) GetMessageInfo(idMess uint64) (utilities.Message, error) {
	var msg utilities.Message
	err := db.c.QueryRow(`SELECT text FROM message WHERE id = ?`, idMess).Scan(&msg.Text)
	if errors.Is(err, sql.ErrNoRows) {
		return utilities.Message{}, ErrMessageNotFound
	}
	return msg, err
}

func (db *appdbimpl) AddMessage(mess *utilities.Message) error {
	// Insert the new message in the database
	err := db.c.QueryRow(`INSERT INTO messages (text, conv_id, sender_id) VALUES (?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Conv, mess.Sender).Scan(&mess.ID, &mess.Timestamp)
	if err != nil {
		return fmt.Errorf("error adding message to database: %v", err)
	}

	// Get the receiver id for insert the status message
	receiver, err := db.GetReceiver(mess.Conv, mess.Sender)
	if err != nil {
		return fmt.Errorf("error getting receiver: %v", err)
	}
	mess.Status, err = db.UpdateStatus(receiver, mess.Sender)
	if err != nil {
		return fmt.Errorf("error updating receiver: %v", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMessage(messId uint64) error {
	_, err := db.c.Exec(`DELETE FROM messages WHERE id = ?`, messId)
	if err != nil {
		return fmt.Errorf("error removing the message from the database: %v", err)
	}
	return nil
}

func (db *appdbimpl) UpdateStatus(receiver uint64, idMess uint64) (string, error) {
	// TODO: make the function globally for any case (read, received, new message)
	var info string
	err := db.c.QueryRow(`INSERT INTO status(receiver_id, mess_id) VALUES (?, ?) RETURNING info`, receiver, idMess).Scan(&info)
	if err != nil {
		return "", fmt.Errorf("error updating status for message to receiver: %v", err)
	}
	return info, nil
}
