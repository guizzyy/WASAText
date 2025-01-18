package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) GetConversations(uID uint64) ([]utilities.Conversation, error) {
	// Get info for the conversation in the homepage (plus unread messages)
	query := `SELECT 
    				c.id,
    				c.type,
    				c.name, 
    				c.photo, 
    				m.id,
    				m.text, 
    				m.photo, 
    				MAX(m.timestamp) as last_message_time, 
    				COUNT(CASE WHEN sm.info = 'Unreceived' THEN m.id END)
				FROM 
				    conversation AS c, message AS m, membership AS ms, status AS sm 
				WHERE 
				    c.id = m.conv_id AND 
				    c.id = ms.conv_id AND 
				    ms.user_id = ? AND 
				    m.id = sm.mess_id AND 
				    sm.receiver_id = ?
				GROUP BY 
				    c.id 
				ORDER BY 
				    last_message_time DESC`
	rows, err := db.c.Query(query, uID, uID)
	if err != nil {
		return nil, fmt.Errorf("error in getting conversations for the homepage: %w", err)
	}
	defer rows.Close()

	// Scan the rows to get conversation info and updating the status
	var convs []utilities.Conversation
	for rows.Next() {
		var mess utilities.Message
		var conv utilities.Conversation
		if err = rows.Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo, &mess.ID, &mess.Text, &mess.Photo, &mess.Timestamp, &conv.CountUnread); err != nil {
			return nil, fmt.Errorf("error in scanning conversations for the homepage: %w", err)
		}
		if err = db.UpdateReceivedStatus(conv.ID, uID); err != nil {
			return nil, fmt.Errorf("error in updating received status: %w", err)
		}
		conv.LastMessage = mess
		convs = append(convs, conv)
	}

	// Check errors during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetConversations: %w", err)
	}
	return convs, nil
}

func (db *appdbimpl) GetConversation(convID uint64, uID uint64) ([]utilities.Message, error) {
	// Get all the messages for the given conversation
	query := `SELECT
    				m.id,
    				m.text,
    				m.photo,
    				m.conv_id,
    				m.sender_id,
    				m.is_forwarded,
    				m.timestamp
				FROM
				    message AS m, status AS s
				WHERE
				    m.id = s.mess_id AND
				    m.conv_id = ?
				ORDER BY m.timestamp DESC`
	rows, err := db.c.Query(query, convID)
	if err != nil {
		return nil, fmt.Errorf("error in getting messages in a conversation: %w", err)
	}
	defer rows.Close()

	// Scan the rows to get messages info and updating the status
	var messages []utilities.Message
	for rows.Next() {
		var m utilities.Message
		pMess := &m
		if err = rows.Scan(&m.ID, &m.Text, &m.Photo, &m.Conv, &m.Sender, &m.IsForward, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("error in scanning messages in a conversation: %w", err)
		}
		if err = db.UpdateReadStatus(pMess, convID, uID); err != nil {
			return nil, fmt.Errorf("error in updating read status: %w", err)
		}
		messages = append(messages, m)
	}

	// Check errors during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in getting messages in a conversation: %w", err)
	}
	return messages, nil
}

func (db *appdbimpl) CreatePrivConv(u utilities.User, receiver utilities.User) (utilities.Conversation, error) {
	// Control if the user and receiver are the same
	if u.ID == receiver.ID {
		return utilities.Conversation{}, errors.New("cannot create a conversation with yourself")
	}

	// Control if the conversation already exists
	if exists, conv, err := db.PrivConvExists(uID, receiver); err != nil {
		return conv, err
	} else if exists {
		return conv, nil
	}

	// Insert the new private conversation created in the database
	var conv utilities.Conversation
	err := db.c.QueryRow(`INSERT INTO conversation(type, name, photo) VALUES (?, ?, ?) RETURNING *`, "private", receiver.Username, receiver.Photo).Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo)
	if err != nil {
		return utilities.Conversation{}, fmt.Errorf("error in creating private conversation: %w", err)
	}

	// Insert the membership of the user in the database
	_, err = db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, conv.ID, uID)
	if err != nil {
		return utilities.Conversation{}, fmt.Errorf("error in creating memberships: %w", err)
	}

	// Insert the membership of the receiver in the database
	_, err = db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, conv.ID, receiver.ID)
	if err != nil {
		return utilities.Conversation{}, fmt.Errorf("error in creating memberships: %w", err)
	}
	return conv, nil
}

func (db *appdbimpl) CreateGroupConv(grConv *utilities.Conversation, user_id uint64) error {
	// Insert and retrieve the new conversation info in the database
	err := db.c.QueryRow(`INSERT INTO conversation(name, type) VALUES (?, ?) RETURNING id, photo`, grConv.Name, grConv.Type).Scan(&grConv.ID, &grConv.Photo)
	if err != nil {
		return fmt.Errorf("error in creating conversation: %w", err)
	}

	// Insert the new membership of the group creator and the new group created
	_, err = db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, grConv.ID, user_id)
	if err != nil {
		return fmt.Errorf("error in adding memberships while creating the group: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetGroupName(group utilities.Conversation, uID uint64) error {
	// Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(group.ID); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return fmt.Errorf("conversation is not a group")
	}

	// Check if the user who is trying to change the name is in the group
	if isIn, err := db.IsUserInConv(group.ID, uID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return fmt.Errorf("user is not in the conversation")
	}

	// Update the group name and check possible errors
	res, err := db.c.Exec(`UPDATE conversation SET name = ? WHERE id = ?`, group.Name, group.ID)
	if err != nil {
		return fmt.Errorf("error in setting group name: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get affected rows in group name db function: %w", err)
	}
	if rows == 0 {
		return ErrConversationNotFound
	}
	return nil
}

func (db *appdbimpl) SetGroupPhoto(group utilities.Conversation, uID uint64) error {
	// Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(group.ID); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return fmt.Errorf("conversation is not a group")
	}

	// Check if the user who is trying to change the photo is in the group
	if isIn, err := db.IsUserInConv(group.ID, uID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return fmt.Errorf("user is not in the conversation")
	}

	// Update the group photo and check possible errors
	res, err := db.c.Exec(`UPDATE conversation SET photo = ? WHERE id = ?`, group.Photo, group.ID)
	if err != nil {
		return fmt.Errorf("error in setting group photo: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get affected rows in set group photo db function: %w", err)
	}
	if rows == 0 {
		return ErrConversationNotFound
	}
	return nil
}

func (db *appdbimpl) AddToGroup(idConv uint64, uID uint64, uAdded utilities.User) error {
	// Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(idConv); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return fmt.Errorf("conversation is not a group")
	}

	// Check if the user is in the group in order to add a new member
	if isIn, err := db.IsUserInConv(idConv, uID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return fmt.Errorf("the user is not in the group conversation")
	}

	// Check if the user added is already in the group
	if isIn, err := db.IsUserInConv(idConv, uAdded.ID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if isIn {
		return fmt.Errorf("the user is already in the group conversation")
	}

	// Insert the new membership of the user to the group conversation
	_, err := db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, idConv, uAdded.ID)
	if err != nil {
		return fmt.Errorf("error in adding membership to conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(idConv uint64, idUser uint64) error {
	// Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(idConv); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return fmt.Errorf("conversation is not a group")
	}

	// Check if the user is in the group
	if isIn, err := db.IsUserInConv(idConv, idUser); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return fmt.Errorf("the user is not in the group conversation")
	}

	// Delete the membership of the user id from the conversation
	_, err := db.c.Exec(`DELETE FROM membership WHERE conv_id = ? AND user_id = ?`, idConv, idUser)
	if err != nil {
		return fmt.Errorf("error in leaving conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetReceivers(convID uint64, senderID uint64) ([]uint64, error) {
	var receivers []uint64
	// Get the set of receivers for a given conversation (if private, it will be an array of 1 element)
	rows, err := db.c.Query(`SELECT user_id FROM membership WHERE conv_id = ? AND user_id != ?`, convID, senderID)
	if err != nil {
		return nil, fmt.Errorf("error in getting receivers of the message: %w", err)
	}
	defer rows.Close()

	// Scan the rows to get the receivers id
	for rows.Next() {
		var receiver uint64
		if err = rows.Scan(&receiver); err != nil {
			return nil, fmt.Errorf("error in scanning receivers of the message: %w", err)
		}
		receivers = append(receivers, receiver)
	}

	// Check error during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetReceivers: %w", err)
	}
	return receivers, nil
}

func (db *appdbimpl) GetMembers(convID uint64, uID uint64) ([]utilities.User, error) {
	// Check if the user is in the conversation in order to use the query
	if isIn, err := db.IsUserInConv(convID, uID); err != nil {
		return nil, fmt.Errorf("error in checking if user is in conversation: %w", err)
	} else if !isIn {
		return nil, fmt.Errorf("user is not in the conversation")
	}

	var members []utilities.User
	// Select infos of user in the conversation
	query := `SELECT 
    				u.id, u.name, u.photo 
				FROM 
				    user AS u, membership AS ms 
				WHERE 
				    ms.user_id = u.id AND ms.conv_id = ?`
	rows, err := db.c.Query(query, convID)
	if err != nil {
		return nil, fmt.Errorf("error in getting members of the conversation: %w", err)
	}
	defer rows.Close()

	// Scan the rows to get id, name and photo of each user
	for rows.Next() {
		var u utilities.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Photo); err != nil {
			return nil, fmt.Errorf("error in scanning members of the conversation: %w", err)
		}
		members = append(members, u)
	}

	// Check error in the resulting rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetMembers: %w", err)
	}
	return members, nil
}

func (db *appdbimpl) GetConvPhoto(convID uint64) (string, error) {
	// Query the database to get the current group photo stored
	var groupPhoto string
	err := db.c.QueryRow(`SELECT photo FROM conversation WHERE id = ?`, convID).Scan(&groupPhoto)
	if err != nil {
		return "", fmt.Errorf("error in getting photo of the conversation: %w", err)
	}
	return groupPhoto, nil
}

func (db *appdbimpl) IsGroupConv(convID uint64) (bool, error) {
	// Check if a given conv id refers to a group conversation
	var class string
	row := db.c.QueryRow(`SELECT type FROM conversation WHERE id = ?`, convID).Scan(&class)
	if row == nil {
		return false, ErrConversationNotFound
	}
	return class == "group", nil
}

func (db *appdbimpl) IsUserConv(convID uint64) (bool, error) {
	// Check if a given conv id refers to a private conversation
	var class string
	row := db.c.QueryRow(`SELECT type FROM conversation WHERE id = ?`, convID).Scan(&class)
	if row == nil {
		return false, ErrConversationNotFound
	}
	return class == "private", nil
}

func (db *appdbimpl) IsUserInConv(convID uint64, uID uint64) (bool, error) {
	// Check if a given user is in the conversation
	_, err := db.c.Query(`SELECT * FROM membership WHERE conv_id = ? AND user_id = ?`, convID, uID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error in checking if a user is in conversation: %w", err)
	}
	return true, nil
}

func (db *appdbimpl) PrivConvExists(u utilities.User, receiver utilities.User) (bool, utilities.Conversation, error) {
	// Select conv information to see if it exists
	var conv utilities.Conversation
	query := `SELECT
					c.id, c,type, c.name, c.photo
				FROM
					conversation AS c,
					membership AS m
				WHERE
				    c.id = m.conv_id AND
				    c.name = ? AND
				    m.user_id = ?`
	err := db.c.QueryRow(query, receiver.Username, uID).Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo)
	if err != nil {
		// If the error is NoRows, it means the conversation doesn't exist
		if errors.Is(err, sql.ErrNoRows) {
			return false, conv, nil
		}
		return false, conv, fmt.Errorf("error in checking if a conversation exists: %w", err)
	}
	return true, conv, nil
}
