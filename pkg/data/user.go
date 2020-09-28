package data

import (
	"github.com/google/uuid"
	"github.com/souhub/wecircles/pkg/logging"
)

type User struct {
	Id        int
	Uuid      string
	Name      string `validate:"required"`
	UserIdStr string `validate:"alphanumunicode"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
	ImgPass   string
	CreatedAt string
}

// Get all users from users
func Users() (users []User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM users`
	rows, err := db.Query(query)
	if err != nil {
		logging.Warn("Failed to find users.")
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImgPass, &user.CreatedAt)
		if err != nil {
			logging.Warn("Failed to find a user.")
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	db := NewDB()
	defer db.Close()
	session = Session{}
	query := `SELECT * FROM sessions
			  WHERE user_id = ?`
	err = db.QueryRow(query, user.Id).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.UserIdStr, &session.CreatedAt)
	return
}

// Get the user from his email address
func UserByEmail(email string) (user User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT * FROM users
			  WHERE email=?`
	err = db.QueryRow(query, email).Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImgPass, &user.CreatedAt)
	return
}

// Get the user from his uuid
// func UserByUserIdStr(user_id_str string) (user User, err error) {
// 	db := NewDB()
// 	defer db.Close()
// 	query := `SELECT * FROM users
// 			  WHERE user_id_str=?`
// 	err = db.QueryRow(query, user_id_str).Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImgPass, &user.CreatedAt)
// 	return user, err
// }
func UserByUserIdStr(useridstr string) (user User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT id, name, user_id_str, email, password, created_at FROM users
			  WHERE user_id_str=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn("Failed to prepare the statements;")
	}
	defer stmt.Close()
	err = stmt.QueryRow(useridstr).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

// Create a new user
func (user *User) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO users (name, user_id_str, email, password)
			  VALUES (?,?,?,?)`
	_, err = db.Exec(query, user.Name, user.UserIdStr, user.Email, user.Password)
	return
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO sessions (uuid, email, user_id, user_id_str)
			  VALUES (?,?,?,?)`
	_, err = db.Exec(query, uuid.New().String(), user.Email, user.Id, user.UserIdStr)
	if err != nil {
		logging.Warn("Failed to insert the new session")
	}
	query = `SELECT * from sessions
			 WHERE user_id=?`
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	err = stmt.QueryRow(user.Id).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.UserIdStr, &session.CreatedAt)
	return
}

// Create a new session for an existing user
// func (user *User) CreateSession() (session Session, err error) {
// 	db := NewDB()
// 	defer db.Close()
// 	query := `INSERT INTO sessions (uuid, email, user_id, user_id_str)
// 			  VALUES (?,?,?,?)`
// 	_, err = db.Exec(query, uuid.New().String(), user.Email, user.Id, user.UserIdStr)
// 	if err != nil {
// 		logging.Warn("Failed to insert the new session")
// 	}
// 	query = `SELECT * from sessions
// 			 WHERE user_id_str=?`
// 	stmt, err := db.Prepare(query)
// 	defer stmt.Close()
// 	err = stmt.QueryRow(user.UserIdStr).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.UserIdStr, &session.CreatedAt)
// 	return
// }

// Update the user
func (user *User) Update() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE users
		      SET name=?, user_id_str=?, email=?,password=?, img_path=?
			  WHERE id=?`
	_, err = db.Exec(query, user.Name, user.UserIdStr, user.Email, user.Password, user.ImgPass, user.Id)
	return
}

// Delete the user
func (user *User) Delete() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from users
			  WHERE id=?`
	_, err = db.Exec(query, user.Id)
	return
}

// Delete all of the users
func ResetUsers() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from users`
	_, err = db.Exec(query)
	return
}
