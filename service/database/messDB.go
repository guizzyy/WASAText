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
	// Insert the new message in the database (also with the case of forward message)
	if mess.IsForward {
		err := db.c.QueryRow(`INSERT INTO messages (text, conv_id, sender_id, is_forwarded) VALUES (?, ?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Conv, mess.Sender, mess.IsForward).Scan(&mess.ID, &mess.Timestamp)
		if err != nil {
			return fmt.Errorf("error adding forwarded message to database: %v", err)
		}
	} else {
		err := db.c.QueryRow(`INSERT INTO messages (text, conv_id, sender_id) VALUES (?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Conv, mess.Sender).Scan(&mess.ID, &mess.Timestamp)
		if err != nil {
			return fmt.Errorf("error adding message to database: %v", err)
		}
	}
	// Get the receivers ids for insert into status message
	receivers, err := db.GetReceivers(mess.Conv, mess.Sender)
	if err != nil {
		return fmt.Errorf("error getting receivers: %v", err)
	}

	// Set the message status in the database (create a new row)
	mess.Status, err = db.InsertStatus(receivers, mess.Sender)
	if err != nil {
		return fmt.Errorf("error updating message status for receivers: %v", err)
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

func (db *appdbimpl) InsertStatus(receivers []uint64, idMess uint64) (string, error) {
	var info string
	// Iterate the array in order to set status for each receiver
	for _, receiver := range receivers {
		err := db.c.QueryRow(`INSERT INTO status(receiver_id, mess_id) VALUES (?, ?) RETURNING info`, receiver, idMess).Scan(&info)
		if err != nil {
			return "", fmt.Errorf("error inserting status for message to receiver: %v", err)
		}
	}
	return info, nil
}
