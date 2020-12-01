package data

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	TwitterID      string
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
	err = stmt.QueryRow(userIdStr).Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.OwnerID, &circle.OwnerIDStr, &circle.TwitterID, &circle.CreatedAt)
	return circle, err
}

// ユーザーがサークルを作ったことがあるかを調べるため（ユーザーIDを変更すると2つめのサークルを作れてしまうバグが発生する）
func GetCirclebyUserID(id int) (circle Circle, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT * FROM circles
			  WHERE owner_id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.OwnerID, &circle.OwnerIDStr, &circle.TwitterID, &circle.CreatedAt)
	return circle, err
}

// Get the owner
func (circle *Circle) GetOwner() (user User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT id, name, user_id_str, image_path
			  FROM users
			  WHERE id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(circle.OwnerID).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.ImagePath)
	return user, err
}

func (circle *Circle) MembershipsByCircleID() (memberships []Membership, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			FROM memberships
			WHERE circle_id=?`
	rows, err := db.Query(query, circle.ID)
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

// Get the number of circle's memberships
func (circle *Circle) CountMemberships() (numberOfMemberships int, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT COUNT(*)
			FROM memberships
			WHERE circle_id=?`
	rows, err := db.Query(query, circle.ID)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	for rows.Next() {
		err = rows.Scan(&numberOfMemberships)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	}
	rows.Close()
	return
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
		err = rows.Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.OwnerID, &circle.OwnerIDStr, &circle.TwitterID, &circle.CreatedAt)
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

// Get the circle from owner_id
func CirclebyOwnerID(id string) (circle Circle, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM circles
			  WHERE owner_id_str=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	err = stmt.QueryRow(id).Scan(&circle.ID, &circle.Name, &circle.ImagePath, &circle.Overview, &circle.Category, &circle.OwnerID, &circle.OwnerIDStr, &circle.TwitterID, &circle.CreatedAt)
	return circle, err
}

func (circle *Circle) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO circles (name, image_path, overview, category, owner_id, owner_id_str, twitter_id)
			  VALUES (?,?,?,?,?,?,?)`
	_, err = db.Exec(query, circle.Name, circle.ImagePath, circle.Overview, circle.Category, circle.OwnerID, circle.OwnerIDStr, circle.TwitterID)
	return
}

func (circle *Circle) Update() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE circles
			  SET name=?,image_path=?, overview=?, category=?, twitter_id=?
			  WHERE id=?`
	_, err = db.Exec(query, circle.Name, circle.ImagePath, circle.Overview, circle.Category, circle.TwitterID, circle.ID)
	return
}

func (circle *Circle) Upload(r *http.Request) (uploadedFileName string, err error) {
	// Allow the "POST" method, only
	if r.Method != "POST" {
		err = errors.New("method error: POST only")
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Make thumbnail dir.
	currentRootDir, err := os.Getwd()
	circleImageDir := fmt.Sprintf("%s/web/img/user%d/circles/mycircle", currentRootDir, circle.OwnerID)
	_, err = os.Stat(circleImageDir)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		if err = os.MkdirAll(circleImageDir, 0777); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
	}
	// Get the file sent form the form
	file, fileHeader, err := r.FormFile("image")
	// Get the uploaded file's name from the file.
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		uploadedFileName = circle.ImagePath
		return
	}
	// Delete the current thumbnail.
	if err = circle.DeleteCircleImage(); err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	uploadedFileName = fileHeader.Filename

	// Set the uploaded file's path
	imagePath := fmt.Sprintf("web/img/user%d/circles/mycircle/%s", circle.OwnerID, uploadedFileName)

	// Save the uploaded file to "imagePath"
	saveImage, err := os.Create(imagePath)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}

	// Write the uploaded file to the file for saving.
	_, err = io.Copy(saveImage, file)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Close the "saveImage" and "file"
	defer saveImage.Close()
	defer file.Close()

	// Upload to S3
	if err = S3Upload(imagePath); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}

	// Delete the post directory on the server
	if err = os.Remove(imagePath); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	return uploadedFileName, err
}

func (circle *Circle) Delete() (err error) {
	return
}

func (circle *Circle) DeleteCircleImage() (err error) {
	// currentRootDir, err := os.Getwd()
	circleImage := fmt.Sprintf("web/img/user%d/circles/mycircle/%s", circle.OwnerID, circle.ImagePath)
	// _, err = os.Stat(circleImage)
	// if err != nil {
	// 	return
	// }
	// err = os.Remove(circleImage)
	if err = S3Delete(circleImage); err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
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
