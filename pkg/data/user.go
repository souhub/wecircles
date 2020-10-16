package data

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/souhub/wecircles/pkg/logging"
)

type User struct {
	Id        int
	Name      string `validate:"required"`
	UserIdStr string `validate:"alphanumunicode"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
	ImagePath string
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
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImagePath, &user.CreatedAt)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// Get all of the posts by the user.
func (user *User) PostsByUser() (posts []Post, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM posts
			  WHERE user_id=?
			  ORDER BY id DESC`
	rows, err := db.Query(query, user.Id)
	if err != nil {
		log.Fatal(err)
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.ThumbnailPath, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		posts = append(posts, post)
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
	err = db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.Email, &user.Password, &user.ImagePath, &user.CreatedAt)
	return
}

// Get the user from user_id_str
func UserByUserIdStr(useridstr string) (user User, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT id, name, user_id_str, password, image_path, created_at FROM users
			  WHERE user_id_str=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	defer stmt.Close()
	err = stmt.QueryRow(useridstr).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.Password, &user.ImagePath, &user.CreatedAt)
	return user, err
}

// Create a new user
func (user *User) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO users (name, user_id_str, email, password, image_path)
			  VALUES (?,?,?,?,?)`
	_, err = db.Exec(query, user.Name, user.UserIdStr, user.Email, user.Password, user.ImagePath)
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
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	query = `SELECT * from sessions
			 WHERE user_id=?`
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	err = stmt.QueryRow(user.Id).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.UserIdStr, &session.CreatedAt)
	return
}

// Update the user
func (user *User) Update() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE users
		      SET name=?, user_id_str=?
			  WHERE id=?`
	_, err = db.Exec(query, user.Name, user.UserIdStr, user.Id)
	return
}

// Update the user image
func (user *User) UpdateImage() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE users
		      SET image_path=?
			  WHERE id=?`
	_, err = db.Exec(query, user.ImagePath, user.Id)
	return
}

// Upload the user's image
func (user *User) Upload(r *http.Request) (uploadedFileName string, err error) {
	// Allow the "POST" method, only
	if r.Method != "POST" {
		err = errors.New("method error: POST only")
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// Parse the form
	err = r.ParseForm()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Get the file sent form the form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Get the uploaded file's name from the file.
	uploadedFileName = fileHeader.Filename
	// Set the uploaded file's path
	imagePath := fmt.Sprintf("web/img/user%d/%s", user.Id, uploadedFileName)

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
	return uploadedFileName, err
}

// Upload the circle's image
func (user *User) UploadCircleImage(r *http.Request) (uploadedFileName string, err error) {
	// Allow the "POST" method, only
	if r.Method != "POST" {
		err = errors.New("method error: POST only")
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// Parse the form
	err = r.ParseForm()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Get the file sent form the form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Get the uploaded file's name from the file.
	uploadedFileName = fileHeader.Filename
	// Set the uploaded file's path
	imagePath := fmt.Sprintf("web/img/user%d/circle/%s", user.Id, uploadedFileName)

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
	return uploadedFileName, err
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

// Delete the user image to update the new one
func (user *User) DeleteUserImage() error {
	currentDir, err := os.Getwd()
	userImage := fmt.Sprintf("%s/web/img/user%d", currentDir, user.Id)
	_, err = os.Stat(userImage)
	if err != nil {
		return err
	}
	err = os.RemoveAll(userImage)
	return err
}

// Delete the user's posts
func (user *User) DeletePosts() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from posts
			  WHERE user_id=?`
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
