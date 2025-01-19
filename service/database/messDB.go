package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) GetMessageInfo(idMess uint64) (utilities.Message, error) {
	// Check if the message exists
	if exists, err := db.IsMessageInDatabase(idMess); err != nil {
		return utilities.Message{}, fmt.Errorf("error checking message existence: %w", err)
	} else if !exists {
		return utilities.Message{}, errors.New(ErrMessageNotFound.Error())
	}

	// Retrieve the information about the message
	var msg utilities.Message
	err := db.c.QueryRow(`SELECT text, photo FROM message WHERE id = ?`, idMess).Scan(&msg.Text, &msg.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utilities.Message{}, ErrMessageNotFound
		}
		return utilities.Message{}, fmt.Errorf("error in getting message info: %w", err)
	}
	return msg, err
}

func (db *appdbimpl) AddMessage(mess *utilities.Message) error {
	// Check if the conversation in which the message is sent exists
	if exists, err := db.IsConvInDatabase(mess.Conv); err != nil {
		return fmt.Errorf("error in checking if the conversation exists: %w", err)
	} else if !exists {
		return fmt.Errorf(ErrConversationNotFound.Error())
	}

	// Check if the user who is trying to send a message is in the conversation
	if isIn, err := db.IsUserInConv(mess.Conv, mess.Sender); err != nil {
		return fmt.Errorf("error in checking if the user is in conversation: %w", err)
	} else if !isIn {
		return fmt.Errorf("user is not in conversation")
	}

	// Insert the new message in the database (also with the case of forward message)
	if mess.IsForward {
		err := db.c.QueryRow(`INSERT INTO message (text, photo, conv_id, sender_id, is_forwarded) VALUES (?, ?, ?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Photo, mess.Conv, mess.Sender, mess.IsForward).Scan(&mess.ID, &mess.Timestamp)
		if err != nil {
			return fmt.Errorf("error adding forwarded message to database: %v", err)
		}
	} else {
		err := db.c.QueryRow(`INSERT INTO message (text, photo, conv_id, sender_id) VALUES (?, ?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Photo, mess.Conv, mess.Sender).Scan(&mess.ID, &mess.Timestamp)
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
	mess.Status, err = db.InsertStatus(receivers, mess.ID, mess.Conv)
	if err != nil {
		return fmt.Errorf("error updating message status for receivers: %v", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMessage(messId uint64, uID uint64) error {
	// Check if the message to delete exists
	if exists, err := db.IsMessageInDatabase(messId); err != nil {
		return fmt.Errorf("error in checking if the message exists: %w", err)
	} else if !exists {
		return fmt.Errorf(ErrMessageNotFound.Error())
	}

	// Check if the message belongs to the owner
	if check, err := db.IsOwnerMessage(messId, uID); err != nil {
		return fmt.Errorf("error checking message owner: %v", err)
	} else if !check {
		return fmt.Errorf("message not owned by this user")
	}

	// Delete the message owned by the user
	_, err := db.c.Exec(`DELETE FROM message WHERE id = ? AND sender_id = ?`, messId, uID)
	if err != nil {
		return fmt.Errorf("error removing the message from the database: %v", err)
	}
	return nil
}

func (db *appdbimpl) InsertStatus(receivers []uint64, idMess uint64, idConv uint64) (string, error) {
	var info string
	// Iterate the array in order to set status for each receiver
	for _, receiver := range receivers {
		err := db.c.QueryRow(`INSERT INTO status(receiver_id, mess_id, conv_id) VALUES (?, ?, ?) RETURNING info`, receiver, idMess, idConv).Scan(&info)
		if err != nil {
			return "", fmt.Errorf("error inserting status for message to receiver: %v", err)
		}
	}
	return info, nil
}

func (db *appdbimpl) UpdateReadStatus(cID uint64, uID uint64) error {
	_, err := db.c.Exec(`UPDATE status SET info = 'Read' WHERE conv_id = ? AND receiver_id = ?`, cID, uID)
	if err != nil {
		return fmt.Errorf("error updating status for read message to receiver: %v", err)
	}
	return nil
}

func (db *appdbimpl) UpdateReceivedStatus(uID uint64) error {
	_, err := db.c.Exec(`UPDATE status SET info = 'Received' WHERE receiver_id = ? AND info = 'Unreceived'`, uID)
	if err != nil {
		return fmt.Errorf("error updating status for received message to receiver: %v", err)
	}
	return nil
}

func (db *appdbimpl) IsOwnerMessage(mID uint64, owner_id uint64) (bool, error) {
	var owner uint64
	err := db.c.QueryRow(`SELECT sender_id FROM message WHERE id = ?`, mID).Scan(&owner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrMessageNotFound
		}
		return false, fmt.Errorf("error in getting owner message: %w", err)
	}
	return owner == owner_id, nil
}

func (db *appdbimpl) IsMessageInConv(mID uint64, cID uint64) (bool, error) {
	// Check if ag
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM message WHERE id = ? AND conv_id = ?)`, mID, cID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error in checking if the message exists in the conversation: %w", err)
	}
	return exists, nil
}
