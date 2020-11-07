package data

import (
	"github.com/souhub/wecircles/pkg/logging"
)

type Membership struct {
	ID       int
	UserID   int
	CircleID int
}

func (membership *Membership) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT memberships (user_id, circle_id)
			  VALUE (?,?)`
	_, err = db.Exec(query, membership.UserID, membership.CircleID)
	return
}

func (user *User) MembershipsByUserID() (memberships []Membership, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			FROM memberships
			WHERE user_id=?`
	rows, err := db.Query(query, user.Id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	for rows.Next() {
		var membership Membership
		err = rows.Scan(&membership.ID, &membership.UserID, &membership.CircleID)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		memberships = append(memberships, membership)
	}
	rows.Close()
	return
}

func GetUsersByUserID(memberships []Membership) (users []User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT id, name, user_id_str, image_path
	FROM users
	WHERE id=?`
	for _, membership := range memberships {
		var user User
		rows, err := db.Query(query, membership.UserID)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		for rows.Next() {
			err = rows.Scan(&user.Id, &user.Name, &user.UserIdStr, &user.ImagePath)
			if err != nil {
				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			}
		}
		users = append(users, user)
	}
	return
}

func (membership *Membership) Circle() (circle Circle, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			FROM circles
			WHERE id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	err = stmt.QueryRow(membership.CircleID).Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.OwnerID, &circle.OwnerIDStr, &circle.TwitterID, &circle.CreatedAt)
	return
}

func (membership *Membership) Check(circle Circle) (valid bool, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM memberships
			  WHERE user_id=?
			  AND circle_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(membership.UserID, circle.ID).Scan(&membership.ID, &membership.UserID, &membership.CircleID)
	if err != nil {
		valid = false
		return
	}
	if membership.ID != 0 {
		valid = true
	}
	return
}

func (membership *Membership) Delete() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE FROM memberships
			WHERE user_id=?
			AND circle_id=?`
	_, err = db.Exec(query, membership.UserID, membership.CircleID)
	return
}
