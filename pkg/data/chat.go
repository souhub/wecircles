package data

import (
	"github.com/souhub/wecircles/pkg/logging"
)

type Chat struct {
	ID               int
	Body             string
	UserID           int
	UserIdStr        string
	UserImagePath    string
	CircleID         int
	CircleOwnerIDStr string
	CreatedAt        string
}

// Get all of the chats
func GetChats(circleID int) (chats []Chat, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM chats
			  WHERE circle_id=?
			  ORDER BY id DESC
			  LIMIT 20`
	rows, err := db.Query(query, circleID)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.Body, &chat.UserID, &chat.UserIdStr, &chat.UserImagePath, &chat.CircleID, &chat.CircleOwnerIDStr, &chat.CreatedAt)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		chats = append(chats, chat)
	}
	defer rows.Close()
	return
}

func GetChatByID(chatID string) (chat Chat, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			FROM chats
			WHERE id=?`
	err = db.QueryRow(query, chatID).Scan(&chat.ID, &chat.Body, &chat.UserID, &chat.UserIdStr, &chat.UserImagePath, &chat.CircleID, &chat.CircleOwnerIDStr, &chat.CreatedAt)
	return
}

func (chat *Chat) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT chats (body, user_id, user_id_str, user_image_path, circle_id, circle_owner_id_str)
			  VALUE (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, chat.Body, chat.UserID, chat.UserIdStr, chat.UserImagePath, chat.CircleID, chat.CircleOwnerIDStr)
	return
}

func (chat *Chat) Delete() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from chats
			WHERE id=?`
	_, err = db.Exec(query, chat.ID)
	return
}
