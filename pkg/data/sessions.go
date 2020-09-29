package data

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	UserIdStr string
	CreatedAt string
}

// Check the session for an existing user
func (session *Session) Check() (valid bool, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT * FROM sessions
			  WHERE uuid=?`
	err = db.QueryRow(query, session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.UserIdStr, &session.CreatedAt)
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
	query := `SELECT id, name, user_id_str, password, image_path,created_at
			  FROM users
			  WHERE id=?`
	err = db.QueryRow(query, session.UserId).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.Password, &user.ImagePath, &user.CreatedAt)
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
