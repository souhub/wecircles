package data

import (
	"time"
)

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	UserIdStr string
	CreatedAt time.Time
}

// Check the session for an existing user
func (session *Session) Check() (valid bool, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT * FROM sessions
			  WHERE uuid=?`
	err = db.QueryRow(query, session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// Get the user from the sessions
func (session *Session) User() (user User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT * FROM users
			  WHERE id=?`
	err = db.QueryRow(query, session.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Password, &user.ImgPass, &user.CreatedAt)
	return
}

// Delete the session
func (session *Session) Delete() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE FROM sessions
			  WHERE uuid=?`
	_, err = db.Exec(query, session.Uuid)
	return err
}

// Delete all of the sessions
func ResetSessions() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from sessions`
	_, err = db.Exec(query)
	return
}
