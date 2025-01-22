package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/utilities"
)

func (db *appdbimpl) GetConversations(uID uint64) ([]utilities.Conversation, error) {
	//	Get info for the conversation in the homepage
	query := `WITH lastMess AS (
				SELECT 
				    conv_id,
				    id AS mess_id,
				    text,
				    photo,
				    timestamp
				FROM
				    message
				WHERE
				    timestamp = (SELECT MAX(timestamp) FROM message AS m2 WHERE m2.conv_id = message.conv_id)
)
				SELECT 
    				c.id,
    				c.type,
    				m.mess_id,
    				m.text, 
    				m.photo, 
    				m.timestamp as last_message_time
				FROM 
				    conversation AS c, lastMess AS m, membership AS ms 
				WHERE 
				    c.id = m.conv_id AND 
				    c.id = ms.conv_id AND 
				    ms.user_id = ? 
				ORDER BY 
				    last_message_time DESC`
	rows, err := db.c.Query(query, uID)
	if err != nil {
		return nil, fmt.Errorf("error in getting conversations for the homepage: %w", err)
	}

	//	Scan the rows to get conversation info and updating the status
	var convs []utilities.Conversation
	for rows.Next() {
		var mess utilities.Message
		var conv utilities.Conversation
		if err = rows.Scan(&conv.ID, &conv.Type, &mess.ID, &mess.Text, &mess.Photo, &mess.Timestamp); err != nil {
			return nil, fmt.Errorf("error in scanning conversations for the homepage: %w", err)
		}

		switch conv.Type {
		case "private":
			if conv.Name, conv.Photo, err = db.GetPrivConvInfo(conv.ID, uID); err != nil {
				return nil, fmt.Errorf("error in getting name, photo of private conversation for the homepage: %w", err)
			}
		case "group":
			if conv.Name, conv.Photo, err = db.GetGroupConvInfo(conv.ID); err != nil {
				return nil, fmt.Errorf("error in getting name, photo of group conversation for the homepage: %w", err)
			}
		}
		mess.Status = "Received"
		conv.LastMessage = mess
		convs = append(convs, conv)
	}

	//	Check errors during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetConversations: %w", err)
	}

	//	Update the status of the messages in the database
	if err = db.UpdateReceivedStatus(uID); err != nil {
		return nil, fmt.Errorf("error in updating received status: %w", err)
	}
	return convs, nil
}

func (db *appdbimpl) GetConversation(convID uint64, uID uint64) ([]utilities.Message, error) {
	//	Check if the conversation is in the database
	if exists, err := db.IsConvInDatabase(convID); err != nil {
		return nil, fmt.Errorf("error checking if conversation is in database: %w", err)
	} else if !exists {
		return nil, ErrConversationNotFound
	}

	//	Check if the user is in the conversation
	if isIn, err := db.IsUserInConv(convID, uID); err != nil {
		return nil, fmt.Errorf("error checking if user is in conversation: %w", err)
	} else if !isIn {
		return nil, ErrUserNotInConversation
	}

	//	Get all the messages for the given conversation
	query := `SELECT
    				m.id,
    				m.text,
    				m.photo,
    				m.conv_id,
    				m.sender_id,
    				m.is_forwarded,
    				m.timestamp
				FROM
				    message AS m
				WHERE
				    m.conv_id = ?
				ORDER BY m.timestamp DESC`
	rows, err := db.c.Query(query, convID)
	if err != nil {
		return nil, fmt.Errorf("error in getting messages in a conversation: %w", err)
	}

	//	Scan the rows to get messages info and updating the status
	var messages []utilities.Message
	for rows.Next() {
		var m utilities.Message
		if err = rows.Scan(&m.ID, &m.Text, &m.Photo, &m.Conv, &m.Sender, &m.IsForward, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("error in scanning messages in a conversation: %w", err)
		}
		m.Status = "Read"
		messages = append(messages, m)
	}

	//	Check errors during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in getting messages in a conversation: %w", err)
	}

	//	Update the status of the messages in the database
	if err = db.UpdateReadStatus(convID, uID); err != nil {
		return nil, fmt.Errorf("error in updating read status: %w", err)
	}
	return messages, nil
}

func (db *appdbimpl) CreatePrivConv(u utilities.User, receiver utilities.User) (utilities.Conversation, error) {
	//	Control if the user and receiver are the same
	if u.ID == receiver.ID {
		return utilities.Conversation{}, ErrNoSelfConversation
	}

	//	Control if the conversation already exists
	if exists, conv, err := db.PrivConvExists(u, receiver); err != nil {
		return conv, err
	} else if exists {
		return conv, nil
	}

	//	Insert the new private conversation created in the database
	var conv utilities.Conversation
	err := db.c.QueryRow(`INSERT INTO conversation(type) VALUES (?) RETURNING id, type`, "private").Scan(&conv.ID, &conv.Type)
	if err != nil {
		return utilities.Conversation{}, fmt.Errorf("error in creating private conversation: %w", err)
	}

	//	Insert the membership of the user in the database
	_, err = db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, conv.ID, u.ID)
	if err != nil {
		return utilities.Conversation{}, fmt.Errorf("error in creating memberships: %w", err)
	}

	// Insert the membership of the receiver in the database
	_, err = db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, conv.ID, receiver.ID)
	if err != nil {
		return utilities.Conversation{}, fmt.Errorf("error in creating memberships: %w", err)
	}

	conv.Name = receiver.Username
	conv.Photo = receiver.Photo
	return conv, nil
}

func (db *appdbimpl) CreateGroupConv(grConv *utilities.Conversation, user_id uint64) error {
	//	Insert and retrieve the new conversation info in the database
	err := db.c.QueryRow(`INSERT INTO conversation(name, type) VALUES (?, ?) RETURNING id`, grConv.Name, grConv.Type).Scan(&grConv.ID)
	if err != nil {
		return fmt.Errorf("error in creating conversation: %w", err)
	}

	//	Insert the new membership of the group creator and the new group created
	_, err = db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, grConv.ID, user_id)
	if err != nil {
		return fmt.Errorf("error in adding memberships while creating the group: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetGroupName(group utilities.Conversation, uID uint64) error {
	//	Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(group.ID); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return ErrNoGroup
	}

	//	Check if the user who is trying to change the name is in the group
	if isIn, err := db.IsUserInConv(group.ID, uID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return ErrUserNotInConversation
	}

	//	Update the group name and check possible errors
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
	//	Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(group.ID); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return ErrNoGroup
	}

	//	Check if the user who is trying to change the photo is in the group
	if isIn, err := db.IsUserInConv(group.ID, uID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return ErrUserNotInConversation
	}

	//	Update the group photo and check possible errors
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
	//	Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(idConv); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return ErrNoGroup
	}

	//	Check if the user is in the group in order to add a new member
	if isIn, err := db.IsUserInConv(idConv, uID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return ErrUserNotInConversation
	}

	//	Check if the user added is already in the group
	if isIn, err := db.IsUserInConv(idConv, uAdded.ID); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if isIn {
		return ErrUserInGroup
	}

	//	Insert the new membership of the user to the group conversation
	_, err := db.c.Exec(`INSERT INTO membership(conv_id, user_id) VALUES (?, ?)`, idConv, uAdded.ID)
	if err != nil {
		return fmt.Errorf("error in adding membership to conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(idConv uint64, idUser uint64) error {
	//	Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(idConv); err != nil {
		return fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return ErrNoGroup
	}

	//	Check if the user is in the group
	if isIn, err := db.IsUserInConv(idConv, idUser); err != nil {
		return fmt.Errorf("error in checking if a user is in the conversation: %w", err)
	} else if !isIn {
		return ErrUserNotInConversation
	}

	//	Delete the membership of the user id from the conversation
	_, err := db.c.Exec(`DELETE FROM membership WHERE conv_id = ? AND user_id = ?`, idConv, idUser)
	if err != nil {
		return fmt.Errorf("error in leaving conversation: %w", err)
	}

	//	Delete the group if no members remain
	if err = db.GroupStillExists(idConv); err != nil {
		return fmt.Errorf("error checking if there are still members in the conversation: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetReceivers(convID uint64, senderID uint64) ([]uint64, error) {
	//	Get the set of receivers for a given conversation (if private, it will be an array of 1 element)
	var receivers []uint64
	rows, err := db.c.Query(`SELECT user_id FROM membership WHERE conv_id = ? AND user_id != ?`, convID, senderID)
	if err != nil {
		return nil, fmt.Errorf("error in getting receivers of the message: %w", err)
	}

	//	Scan the rows to get the receivers id
	for rows.Next() {
		var receiver uint64
		if err = rows.Scan(&receiver); err != nil {
			return nil, fmt.Errorf("error in scanning receivers of the message: %w", err)
		}
		receivers = append(receivers, receiver)
	}

	//	Check error during the scanning of the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetReceivers: %w", err)
	}
	return receivers, nil
}

func (db *appdbimpl) GetMembers(convID uint64, uID uint64) ([]utilities.User, error) {
	//	Check if the conversation is in the database
	if exists, err := db.IsConvInDatabase(convID); err != nil {
		return nil, fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !exists {
		return nil, ErrConversationNotFound
	}

	//	Check if the user is in the conversation in order to use the query
	if isIn, err := db.IsUserInConv(convID, uID); err != nil {
		return nil, fmt.Errorf("error in checking if user is in conversation: %w", err)
	} else if !isIn {
		return nil, ErrUserNotInConversation
	}

	//	Select infos of user in the conversation
	var members []utilities.User
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

	//	Scan the rows to get id, name and photo of each user
	for rows.Next() {
		var u utilities.User
		var photo sql.NullString
		if err = rows.Scan(&u.ID, &u.Username, &photo); err != nil {
			return nil, fmt.Errorf("error in scanning members of the conversation: %w", err)
		}
		u.Photo = photo.String
		members = append(members, u)
	}

	//	Check error in the resulting rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in resulting rows of GetMembers: %w", err)
	}
	return members, nil
}

func (db *appdbimpl) GetConvPhoto(convID uint64) (string, error) {
	//	Check if the id refers to a group conversation
	if isGroup, err := db.IsGroupConv(convID); err != nil {
		return "", fmt.Errorf("error in checking if conversation is a group: %w", err)
	} else if !isGroup {
		return "", ErrNoGroup
	}

	//	Query the database to get the current group photo stored
	var groupPhoto sql.NullString
	err := db.c.QueryRow(`SELECT photo FROM conversation WHERE id = ?`, convID).Scan(&groupPhoto)
	if err != nil {
		return "", fmt.Errorf("error in getting photo of the conversation: %w", err)
	}
	return groupPhoto.String, nil
}

func (db *appdbimpl) IsGroupConv(convID uint64) (bool, error) {
	//	Check if a given conv id refers to a group conversation
	var class string
	err := db.c.QueryRow(`SELECT type FROM conversation WHERE id = ?`, convID).Scan(&class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrConversationNotFound
		}
		return false, fmt.Errorf("error checking if a conversation is a group: %w", err)
	}
	return class == "group", nil
}

func (db *appdbimpl) IsPrivConv(convID uint64) (bool, error) {
	//	Check if a given conv id refers to a private conversation
	var class string
	err := db.c.QueryRow(`SELECT type FROM conversation WHERE id = ?`, convID).Scan(&class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrConversationNotFound
		}
		return false, fmt.Errorf("error checking if a conversation is private: %w", err)
	}
	return class == "private", nil
}

func (db *appdbimpl) IsUserInConv(convID uint64, uID uint64) (bool, error) {
	//	Check if a given user is in the conversation
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM membership WHERE conv_id = ? AND user_id = ?)`, convID, uID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if a user is in the conversation: %w", err)
	} else if !exists {
		return false, nil
	}
	return true, nil
}

func (db *appdbimpl) PrivConvExists(u utilities.User, receiver utilities.User) (bool, utilities.Conversation, error) {
	//	Select conv information to see if it exists
	var conv utilities.Conversation
	query := `SELECT
					c.id, c.type
				FROM
					conversation AS c
				WHERE	
				    c.id IN (
				    	SELECT 
				    	    m.conv_id 
				    	FROM 
				    	    membership AS m 
				    	WHERE 
				    	    m.user_id IN (?, ?) 
				    	GROUP BY 
				    	    m.conv_id 
				    	HAVING 
				    	    COUNT(DISTINCT m.user_id) = 2)`
	err := db.c.QueryRow(query, u.ID, receiver.ID).Scan(&conv.ID, &conv.Type)
	if err != nil {
		//	If the error is NoRows, it means the conversation doesn't exist
		if errors.Is(err, sql.ErrNoRows) {
			return false, conv, nil
		}
		return false, conv, fmt.Errorf("error in checking if a conversation exists: %w", err)
	}
	conv.Name = receiver.Username
	conv.Photo = receiver.Photo
	return true, conv, nil
}

func (db *appdbimpl) GroupStillExists(idConv uint64) error {
	//	Check if there is at least 1 member in the group conversation (otherwise, delete the group)
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM membership WHERE conv_id = ?)`, idConv).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if there is a member in the group conversation: %w", err)
	}
	if !exists {
		if _, err = db.c.Exec(`DELETE FROM conversation WHERE id = ?`, idConv); err != nil {
			return fmt.Errorf("error deleting a group conversation: %w", err)
		}
	}
	return nil
}

func (db *appdbimpl) GetPrivConvInfo(convID uint64, senderID uint64) (string, string, error) {
	//	Get username and photo of the receiver to display in the homepage
	receiver, err := db.GetReceivers(convID, senderID)
	if err != nil {
		return "", "", fmt.Errorf("error in getting receivers of the conversation: %w", err)
	}
	userRec, err := db.GetUserByID(receiver[0])
	if err != nil {
		return "", "", fmt.Errorf("error in getting receiver of the private conversation: %w", err)
	}
	return userRec.Username, userRec.Photo, nil
}

func (db *appdbimpl) GetGroupConvInfo(convID uint64) (string, string, error) {
	//	Get name and photo of the group to display in the homepage
	var group utilities.Conversation
	err := db.c.QueryRow(`SELECT name, photo FROM conversation WHERE id = ?`, convID).Scan(&group.Name, &group.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrConversationNotFound
		}
		return "", "", fmt.Errorf("error in getting group conversation info: %w", err)
	}
	return group.Name, group.Photo, nil
}
