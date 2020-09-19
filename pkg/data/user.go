package data

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/souhub/wecircles/pkg/config"
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
	db := config.NewDB()
	err = db.Table("users").Find(&users).Error
	return
}

func UserByEmail(email string) (user User) {
	db := config.NewDB()
	db.Table("users").Where("email=?", email).Scan(&user)
	return user
}

func UserByUserIdStr(user_id_str string) (user User, err error) {
	db := config.NewDB()
	err = db.Table("users").Where("user_id_str=?", user_id_str).Scan(&user).Error
	return
}

func (user *User) CreateSession() Session {
	session := Session{
		Uuid:      uuid.New().String(),
		Email:     user.Email,
		UserId:    user.Id,
		CreatedAt: time.Now(),
	}
	db := config.NewDB()
	db.Table("sessions").Create(&session)
	return session
}

func (session *Session) Check() (valid bool, err error) {
	db := config.NewDB()
	err = db.Table("sessions").Where("uuid=?", session.Uuid).Scan(&session).Error
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
	user = User{}
	db := config.NewDB()
	err = db.Table("users").Where("id=?", session.UserId).Scan(&user).Error
	return
}

func (session *Session) UpdateUser(userPtr *User, attr map[string]interface{}) (err error) {
	db := config.NewDB()
	err = db.Model(userPtr).Update(attr).Error
	return
}

func (session *Session) DeleteUser(userPtr *User) (err error) {
	db := config.NewDB()
	err = db.Delete(userPtr).Error
	return
}
