package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
	"time"
)

func (db *appdbimpl) GetMessageInfo(idMess uint64) (utilities.Message, error) {
	//	Check if the message exists
	if exists, err := db.IsMessageInDatabase(idMess); err != nil {
		return utilities.Message{}, fmt.Errorf("error checking message existence: %w", err)
	} else if !exists {
		return utilities.Message{}, ErrMessageNotFound
	}

	//	Retrieve the information about the message
	var msg utilities.Message
	var sender uint64
	err := db.c.QueryRow(`SELECT text, photo, sender_id FROM message WHERE id = ?`, idMess).Scan(&msg.Text, &msg.Photo, &sender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utilities.Message{}, ErrMessageNotFound
		}
		return utilities.Message{}, fmt.Errorf("error in getting message info: %w", err)
	}
	if msg.Sender, err = db.GetUserByID(sender); err != nil {
		return utilities.Message{}, fmt.Errorf("error in getting message info: %w", err)
	}
	return msg, err
}

func (db *appdbimpl) AddMessage(mess *utilities.Message) error {
	//	Check if the conversation in which the message is sent exists
	if exists, err := db.IsConvInDatabase(mess.Conv); err != nil {
		return fmt.Errorf("error in checking if the conversation exists: %w", err)
	} else if !exists {
		return ErrConversationNotFound
	}

	//	Check if the user who is trying to send a message is in the conversation
	if isIn, err := db.IsUserInConv(mess.Conv, mess.Sender.ID); err != nil {
		return fmt.Errorf("error in checking if the user is in conversation: %w", err)
	} else if !isIn {
		return ErrUserNotInConversation
	}

	//	Insert the new message in the database (also with the case of forward message)
	if mess.IsForward {
		err := db.c.QueryRow(`INSERT INTO message (text, photo, conv_id, sender_id, is_forwarded) VALUES (?, ?, ?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Photo, mess.Conv, mess.Sender.ID, mess.IsForward).Scan(&mess.ID, &mess.Timestamp)
		if err != nil {
			return fmt.Errorf("error adding forwarded message to database: %w", err)
		}
	} else {
		err := db.c.QueryRow(`INSERT INTO message (text, photo, conv_id, sender_id) VALUES (?, ?, ?, ?) RETURNING id, timestamp`, mess.Text, mess.Photo, mess.Conv, mess.Sender.ID).Scan(&mess.ID, &mess.Timestamp)
		if err != nil {
			return fmt.Errorf("error adding message to database: %w", err)
		}
	}

	//	Get the receivers ids for insert into status message
	receivers, err := db.GetReceivers(mess.Conv, mess.Sender.ID)
	if err != nil {
		return fmt.Errorf("error getting receivers: %w", err)
	}

	//	Set the message status in the database (create a new row)
	mess.Status, err = db.InsertStatus(receivers, mess.ID, mess.Conv)
	if err != nil {
		return fmt.Errorf("error inserting message status for receivers: %w", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMessage(messId uint64, uID uint64) error {
	//	Check if the message to delete exists
	if exists, err := db.IsMessageInDatabase(messId); err != nil {
		return fmt.Errorf("error in checking if the message exists: %w", err)
	} else if !exists {
		return ErrMessageNotFound
	}

	//	Check if the message belongs to the owner
	if check, err := db.IsOwnerMessage(messId, uID); err != nil {
		return fmt.Errorf("error checking message owner: %w", err)
	} else if !check {
		return fmt.Errorf("message not owned by this user")
	}

	//	Delete the message owned by the user
	_, err := db.c.Exec(`DELETE FROM message WHERE id = ? AND sender_id = ?`, messId, uID)
	if err != nil {
		return fmt.Errorf("error removing the message from the database: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetLastMessage(conv utilities.Conversation, uID uint64) (utilities.Message, error) {
	if isIn, err := db.IsConvInDatabase(conv.ID); err != nil {
		return utilities.Message{}, fmt.Errorf("error in checking if the conversation exists: %w", err)
	} else if !isIn {
		return utilities.Message{}, ErrConversationNotFound
	}

	var joinTimestamp time.Time
	if conv.Type == "group" {
		err := db.c.QueryRow(`SELECT timestamp FROM membership WHERE conv_id = ? AND user_id = ?`, conv.ID, uID).Scan(&joinTimestamp)
		if err != nil {
			return utilities.Message{}, fmt.Errorf("error getting conversation membership timestamp: %w", err)
		}
	}

	var msg utilities.Message
	var msgPhoto sql.NullString
	var senderId uint64
	query := `SELECT
    				m.id,
					m.text,
					m.photo,
					m.conv_id,
					m.sender_id,
					m.timestamp
    			FROM 
    			    message AS m 
    			WHERE
    			    m.conv_id = ? AND
    			    m.timestamp >= ?
				ORDER BY 
				    m.timestamp DESC
				LIMIT 1`
	err := db.c.QueryRow(query, conv.ID, joinTimestamp).Scan(&msg.ID, &msg.Text, &msgPhoto, &msg.Conv, &senderId, &msg.Timestamp)
	if errors.Is(err, sql.ErrNoRows) {
		return utilities.Message{}, nil
	} else if err != nil {
		return utilities.Message{}, fmt.Errorf("error in getting last message info: %w", err)
	}
	if msg.Sender, err = db.GetUserByID(senderId); err != nil {
		return utilities.Message{}, fmt.Errorf("error getting sender info: %w", err)
	}
	msg.Photo = msgPhoto.String
	return msg, nil
}

func (db *appdbimpl) InsertStatus(receivers []uint64, idMess uint64, idConv uint64) (string, error) {
	//	Iterate the array in order to set status for each receiver
	var info string
	for _, receiver := range receivers {
		err := db.c.QueryRow(`INSERT INTO status(receiver_id, mess_id, conv_id) VALUES (?, ?, ?) RETURNING info`, receiver, idMess, idConv).Scan(&info)
		if err != nil {
			return "", fmt.Errorf("error inserting status for message to receiver: %w", err)
		}
	}
	return info, nil
}

func (db *appdbimpl) UpdateReadStatus(cID uint64, uID uint64) error {
	//	Update the status of the message to READ
	_, err := db.c.Exec(`UPDATE status SET info = 'Read' WHERE conv_id = ? AND receiver_id = ?`, cID, uID)
	if err != nil {
		return fmt.Errorf("error updating status for read message to receiver: %w", err)
	}
	return nil
}

func (db *appdbimpl) UpdateReceivedStatus(uID uint64) error {
	//	Update the status of the message to RECEIVED
	_, err := db.c.Exec(`UPDATE status SET info = 'Received' WHERE receiver_id = ? AND info = 'Unreceived'`, uID)
	if err != nil {
		return fmt.Errorf("error updating status for received message to receiver: %w", err)
	}
	return nil
}

func (db *appdbimpl) CheckStatus(mID uint64, sID uint64) (string, error) {
	var nReceivers int
	queryRec := `SELECT COUNT(*) FROM membership AS ms
                	WHERE ms.conv_id = (SELECT conv_id FROM message WHERE id = ?) AND
                	      ms.timestamp <= (SELECT timestamp FROM message WHERE id = ?) AND
                	      ms.user_id != ?`
	err := db.c.QueryRow(queryRec, mID, mID, sID).Scan(&nReceivers)
	if err != nil {
		return "", fmt.Errorf("error getting receivers for a message: %w", err)
	}

	// Check the current status of the message for all the users
	query := `SELECT
					s.info, COUNT(*)
				FROM
				    status AS s
				WHERE
				    s.mess_id = ?
				GROUP BY
				    s.info
				`
	rows, err := db.c.Query(query, mID)
	if err != nil {
		return "", fmt.Errorf("error checking status for message to receiver: %w", err)
	}
	defer rows.Close()

	statusCount := map[string]int{
		"Unreceived": 0,
		"Received":   0,
		"Read":       0,
	}

	for rows.Next() {
		var info string
		var count int
		if err = rows.Scan(&info, &count); err != nil {
			return "", fmt.Errorf("error scanning message to receiver: %w", err)
		}
		statusCount[info] = count
	}

	// Handle missing status entries by assuming them as "Unreceived"
	statusCount["Unreceived"] += nReceivers - (statusCount["Read"] + statusCount["Received"])

	//	Check errors during the scanning of the rows
	if err = rows.Err(); err != nil {
		return "", fmt.Errorf("error in getting status in a conversation: %w", err)
	}

	if statusCount["Read"] >= nReceivers {
		return "Read", nil
	}
	if statusCount["Unreceived"] > 0 {
		return "Unreceived", nil
	}
	return "Received", nil
}

func (db *appdbimpl) IsOwnerMessage(mID uint64, owner_id uint64) (bool, error) {
	//	Check if the owner id is the same of the message
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
	//	Check if a message is in the conversation
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
