package data

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/logging"
)

type Post struct {
	Id            int
	Uuid          string
	Title         string
	Body          string
	UserId        int
	UserIdStr     string
	UserName      string
	ThumbnailPath string
	CreatedAt     string
}

// Get all of the posts
func Posts() (posts []Post, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM posts
			  ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.ThumbnailPath, &post.CreatedAt)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		posts = append(posts, post)
	}
	defer rows.Close()
	return
}

// Get the post from uuid
func PostByUuid(uuid string) (post Post, err error) {
	db := NewDB()
	defer db.Close()
	post = Post{}
	query := `SELECT *
			  FROM posts
			  WHERE uuid=?`
	err = db.QueryRow(query, uuid).Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.ThumbnailPath, &post.CreatedAt)
	return
}

// Get the post from uuid
func PostByID(id string) (post Post, err error) {
	db := NewDB()
	defer db.Close()
	post = Post{}
	query := `SELECT *
			  FROM posts
			  WHERE id=?`
	err = db.QueryRow(query, id).Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.ThumbnailPath, &post.CreatedAt)
	return
}

//Get the user from his post
func (post *Post) UserByPost() (user User, err error) {
	db := NewDB()
	defer db.Close()
	user = User{}
	query := `SELECT id, name, user_id_str, email, created_at
			  FROM users
			  WHERE user_id=?`
	err = db.QueryRow(query, post.UserId).Scan(&user.Id, &user.Name, &user.UserIdStr, &user.Email, &user.CreatedAt)
	return
}

// Create a post
func (post *Post) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO posts (uuid, title, body, user_id, user_id_str, user_name, thumbnail_path ,created_at)
			  VALUES (?,?,?,?,?,?,?,?)`
	_, err = db.Exec(query, post.Uuid, post.Title, post.Body, post.UserId, post.UserIdStr, post.UserName, "default_thumbnail.jpg", post.CreatedAt)
	return
}

// Update the post
func (post *Post) UpdatePost() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE posts
			  SET title=?,body=?, thumbnail_path=?
			  WHERE id=?`
	_, err = db.Exec(query, post.Title, post.Body, post.ThumbnailPath, post.Id)
	return
}

// Update the user image
func (post *Post) UpdateThumbnail() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE posts
		      SET thumbnail_path=?
			  WHERE id=?`
	_, err = db.Exec(query, post.ThumbnailPath, post.Id)
	return
}

func (post *Post) UploadThumbnail(r *http.Request) (uploadedFileName string, err error) {
	// Allow the "POST" method, only
	if r.Method != "POST" {
		err = errors.New("method error: POST only")
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}

	// // Delete the current thumbnail.
	// if err = post.DeleteThembnail(); err != nil {
	// 	logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	// }

	// Make thumbnail dir.
	currentRootDir, err := os.Getwd()
	thumbnailImageDir := fmt.Sprintf("%s/web/img/user%d/posts/post%d", currentRootDir, post.UserId, post.Id)
	_, err = os.Stat(thumbnailImageDir)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		if err = os.MkdirAll(thumbnailImageDir, 0777); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
	}
	// Parse the form
	// err = r.ParseForm()
	// if err != nil {
	// 	logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	// 	return
	// }
	// Get the file sent form the form
	file, fileHeader, err := r.FormFile("image")
	// Get the uploaded file's name from the file.
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		uploadedFileName = post.ThumbnailPath
		return
	}
	// Delete the current thumbnail.
	if err = post.DeleteThembnail(); err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	uploadedFileName = fileHeader.Filename

	// Set the uploaded file's path
	imagePath := fmt.Sprintf("web/img/user%d/posts/post%d/%s", post.UserId, post.Id, uploadedFileName)

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

// Delete the thumbnail
func (post *Post) DeleteThembnail() (err error) {
	currentRootDir, err := os.Getwd()
	thumbnail := fmt.Sprintf("%s/web/img/user%d/posts/post%d/%s", currentRootDir, post.UserId, post.Id, post.ThumbnailPath)
	_, err = os.Stat(thumbnail)
	if err != nil {
		return
	}
	err = os.Remove(thumbnail)
	return
}

// Delete the post.
func (post *Post) Delete() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from posts
			  WHERE id=?`
	_, err = db.Exec(query, post.Id)
	return
}

// Delete all oh the posts
func ResetPosts() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from posts`
	_, err = db.Exec(query)
	return
}
