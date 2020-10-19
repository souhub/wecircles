package data

import (
	"log"

	"github.com/souhub/wecircles/pkg/logging"
)

type Circle struct {
	ID             int
	Name           string
	ImagePath      string
	Overview       string
	Category       string
	OwnerID        int
	OwnerIDStr     string
	OwnerImagePath string
	CreatedAt      string
	Owner          User
	Members        []User
}

// Get the owner's circle
func GetCirclebyUser(userIdStr string) (circle Circle, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT * FROM circles
			  WHERE owner_id_str=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(userIdStr).Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.OwnerID, &circle.OwnerIDStr, &circle.CreatedAt)
	return circle, err
}

// Get all of the circles
func Circles() (circles []Circle, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM circles
			  ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	for rows.Next() {
		var circle Circle
		err = rows.Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.Owner.Id, &circle.OwnerIDStr, &circle.CreatedAt)
		if err != nil {
			log.Fatal(err)
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		circles = append(circles, circle)
	}
	defer rows.Close()
	return
}

func (circle *Circle) UserByOwnerID() (user User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT id, name, user_id_str, image_path, created_at
			  FROM users
			  WHERE id=?`
	err = db.QueryRow(query, circle.OwnerID).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.ImagePath, &user.CreatedAt)
	return
}

func (circle *Circle) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO circles (name, image_path, overview, category, owner_id, owner_id_str)
			  VALUES (?,?,?,?,?,?)`
	_, err = db.Exec(query, circle.Name, circle.ImagePath, circle.Overview, circle.Category, circle.OwnerID, circle.OwnerIDStr)
	return
}

func (circle *Circle) Update() (err error) {
	return
}

func (circle *Circle) Delete() (err error) {
	return
}

// Delete all of the users
func ResetCircles() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from circles`
	_, err = db.Exec(query)
	return
}
