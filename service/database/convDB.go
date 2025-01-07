package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) GetConversations(uID uint64) ([]utilities.Conversation, error) {
	// TODO: manage the status of received message

	// Select the conv infos where the user participates ordered by last message
	query := `SELECT 
    			  c.id, c.type, c.name, c.photo, m1.text, m1.timestamp
			  FROM 
			      message AS m1
			      INNER JOIN conversation AS c ON c.id = m1.conv_id  
			      INNER JOIN memberships AS ms ON ms.conv_id = c.id
			  WHERE 
			      ms.user_id = ? AND m1.timestamp = (SELECT MAX(m2.timestamp) FROM message AS m2 WHERE m2.conv_id = m1.conv_id)
			  ORDER BY 
			      m1.timestamp DESC`
	rows, err := db.c.Query(query, uID)
	if err != nil {
		return nil, fmt.Errorf("error in getting conversations info: %w", err)
	}
	defer rows.Close()

	// Create an array of conversation structs to return and scan the rows
	convs := make([]utilities.Conversation, 0)
	for rows.Next() {
		var conv utilities.Conversation
		if err := rows.Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo, &conv.LastMess, &conv.Timestamp); err != nil {
			return nil, fmt.Errorf("error in scanning conversation info: %w", err)
		}
		convs = append(convs, conv)
	}

	// Check errors during the scan, otherwise return the array
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}
	return convs, nil
}

func (db *appdbimpl) GetConversation(convID uint64) ([]utilities.Message, error) {
	// TODO: find a way to manage the status of the message situation
	query := `SELECT 
    			id, text, sender_id, timestamp, info 
			  FROM 
			    message
			  	INNER JOIN status ON status.mess_id = message.id
			  WHERE 
			    conv_id = ?`
	rows, err := db.c.Query(query, convID)
	if err != nil {
		return nil, fmt.Errorf("error in getting conversation: %w", err)
	}
	defer rows.Close()

	messages := make([]utilities.Message, 0)
	for rows.Next() {
		var msg utilities.Message
		msg.Conv = convID
		if err := rows.Scan(&msg.ID, &msg.Text, &msg.Sender, &msg.Timestamp, &msg.Status); err != nil {
			return nil, fmt.Errorf("error in scanning conversation info: %w", err)
		}
		if msg.Status == "Received" {
			msg.Status = "Read"
		}
		messages = append(messages, msg)
	}
}

func (db *appdbimpl) CreateGroupConv(grConv *utilities.Conversation, user_id uint64) error {
	// Insert and retrieve the new conversation info in the database
	err := db.c.QueryRow(`INSERT INTO conversation(name, type) VALUES (?, ?) RETURNING id, photo`, grConv.Name, grConv.Type).Scan(&grConv.ID, &grConv.Photo)
	if err != nil {
		return fmt.Errorf("error in creating conversation: %w", err)
	}

	// Insert the new membership of the group creator and the new group created
	_, err = db.c.Exec(`INSERT INTO memberships(conv_id, user_id) VALUES (?, ?)`, grConv.ID, user_id)
	if err != nil {
		return fmt.Errorf("error in adding memberships while creating the group: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetGroupName(group utilities.Conversation) error {
	res, err := db.c.Exec(`UPDATE conversation SET name = ? WHERE id = ?`, group.Name, group.ID)
	if err != nil {
		return fmt.Errorf("error in setting group name: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get affected rows in group name db function: %w", err)
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (db *appdbimpl) SetGroupPhoto(group utilities.Conversation) error {
	res, err := db.c.Exec(`UPDATE conversation SET photo = ? WHERE id = ?`, group.Photo, group.ID)
	if err != nil {
		return fmt.Errorf("error in setting group photo: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get affected rows in set group photo db function: %w", err)
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (db *appdbimpl) AddToGroup(idConv uint64, u utilities.User) error {
	_, err := db.c.Exec(`INSERT INTO memberships(conv_id, user_id) VALUES (?, ?)`, idConv, u.ID)
	if err != nil {
		return fmt.Errorf("error in adding membership to conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(idConv uint64, idUser uint64) error {
	_, err := db.c.Exec(`DELETE FROM memberships WHERE conv_id = ? AND user_id = ?`, idConv, idUser)
	if err != nil {
		return fmt.Errorf("error in leaving conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetReceivers(convID uint64, senderID uint64) ([]uint64, error) {
	var receivers []uint64
	// Get the set of receivers for a given conversation (if private, it will be an array of 1 element)
	rows, err := db.c.Query(`SELECT user_id FROM memberships WHERE conv_id = ? AND user_id != ?`, convID, senderID)
	if err != nil {
		return nil, fmt.Errorf("error in getting receivers of the message: %w", err)
	}
	defer rows.Close()

	// Scan the rows to get the receivers id
	for rows.Next() {
		var receiver uint64
		if err = rows.Scan(&receiver); err != nil {
			return nil, fmt.Errorf("error in getting receivers of the message: %w", err)
		}
		receivers = append(receivers, receiver)
	}

	// Check error during the scanning of the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetReceivers: %w", err)
	}
	return receivers, nil
}
