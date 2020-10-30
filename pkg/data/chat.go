package data

import (
	"github.com/souhub/wecircles/pkg/logging"
)

type Chat struct {
	ID            int
	Body          string
	UserID        int
	UserIdStr     string
	UserImagePath string
	CircleID      int
	CreatedAt     string
}

// Get all of the chats
func GetChats(circleID int) (chats []Chat, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM chats
			  WHERE circle_id=?
			  ORDER BY id DESC`
	rows, err := db.Query(query, circleID)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.Body, &chat.UserID, &chat.UserIdStr, &chat.UserImagePath, &chat.CircleID, &chat.CreatedAt)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		chats = append(chats, chat)
	}
	defer rows.Close()
	return
}

func (chat *Chat) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT chats (body, user_id, user_id_str, user_image_path, circle_id)
			  VALUE (?, ?, ?, ?, ?)`
	_, err = db.Exec(query, chat.Body, chat.UserID, chat.UserIdStr, chat.UserImagePath, chat.CircleID)
	return
}
