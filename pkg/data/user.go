package data

import (
	"time"

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
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	UserIdStr string
	CreatedAt time.Time
}

func Users() (users []User, err error) {
	defer db.Close()
	cmd := "SELECT * FROM users"
	rows, err := db.Query(cmd)
	if err != nil {
		logging.Warn("Coudn't find users.")
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImgPass, &user.CreatedAt)
		if err != nil {
			logging.Warn("Couldn't find a user.")
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

func UserByEmail(email string) (user User, err error) {
	defer db.Close()
	cmd := "SELECT * FROM users WHERE email=$1"
	err = db.QueryRow(cmd, email).Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImgPass, &user.CreatedAt)
	return
}

func UserByUserIdStr(user_id_str string) (user User, err error) {
	defer db.Close()
	cmd := "SELECT * FROM users WHERE user_id_str=$1"
	err = db.QueryRow(cmd, user_id_str).Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImgPass, &user.CreatedAt)
	return
}

func (user *User) CreateSession() (session Session, err error) {
	statement := "INSERT INTO sessions (uuid,email,user_id,user_id_str) VALUES ($1,$2,$3,$4) RETURNING id,uuid,email,user_id,user_id_str,created_at"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(uuid.New().String(), user.Email, user.Id).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.UserIdStr, &session.CreatedAt)
	return
}

func (session *Session) Check() (valid bool, err error) {
	defer db.Close()
	cmd := "SELECT * FROM sessions WHERE uuid=$1"
	err = db.QueryRow(cmd, session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

func (session *Session) User() (user User, err error) {
	defer db.Close()
	cmd := "SELECT * FROM users WHERE id=$1"
	err = db.QueryRow(cmd, session.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.UserIdStr, &user.Password, &user.ImgPass, &user.CreatedAt)
	return
}

func (user *User) UpdateUser() (err error) {
	defer db.Close()
	statement := "UPDATE users SET name=$2,userIdStr=$3,password=$4,image_path=$5, WHERE id=$1"
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	_, err = stmt.Exec(user.Id, user.Name, user.UserIdStr, user.Password, user.ImgPass)
	return
}

func (session *Session) DeleteUser(user User) (err error) {
	defer db.Close()
	cmd := "DELETE FROM users WHERE id=$1"
	_, err = db.Exec(cmd, user.Id)
	return err
}
