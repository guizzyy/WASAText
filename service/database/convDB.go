package database

func (db *appdbimpl) GetConversations(uID uint64) ([]uint64, error) {
	rows, err := db.c.Query("SELECT id FROM conversations WHERE user1_id = ? OR user2_id = ? UNION SELECT id FROM conversations JOIN membership on conversations.group_id = membership.group_id WHERE membership.user_id = ?", uID, uID, uID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convIDs []uint64
	for rows.Next() {
		var convID uint64
		if err := rows.Scan(&convID); err != nil {
			return nil, err
		}
		convIDs = append(convIDs, convID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return convIDs, nil
}

func (db *appdbimpl) GetConversation(convID uint64) (uint64, error) {
	return 0, nil
}

func (db *appdbimpl) GetConvInfos(convIDs []uint64) ([]uint64, error) {
	return nil, nil
}

/*
func (db *appdbimpl) GetReceiver(convID uint64) (uint64, error) {
	var receiver uint64
	err := db.c.QueryRow(`SELECT user2_id FROM conversations WHERE id = ?`, convID).Scan(&receiver)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("conversation not found")
	}
	return receiver, err
}
*/
