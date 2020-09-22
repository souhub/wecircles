package data

import (
	"github.com/google/uuid"
	"github.com/souhub/wecircles/pkg/logging"
)

type Post struct {
	Id        int
	Uuid      string
	Title     string `validate:"required"`
	Body      string `validate:"required"`
	UserId    int
	UserIdStr string
	UserName  string
	CreatedAt string
}

// Get all of the posts
func Posts() (posts []Post, err error) {
	db := NewDB()
	defer db.Close()
	query := `SELECT *
			  FROM posts`
	rows, err := db.Query(query)
	if err != nil {
		logging.Warn("Failed to find posts.")
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt, &post.CreatedAt)
		if err != nil {
			logging.Warn("Failed to find a post.")
		}
		posts = append(posts, post)
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
			  WHERE user_id=?`
	rows, err := db.Query(query, user.Id)
	if err != nil {
		logging.Warn("Failed to find posts.")
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt, &post.CreatedAt)
		if err != nil {
			logging.Warn("Failed to find a post.")
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Get the post from uuid
func PostByUuid(uuid string) (post Post, err error) {
	db := NewDB()
	defer db.Close()
	post = Post{}
	query := `SELECT *
			  FROM users
			  WHERE uuid=?`
	err = db.QueryRow(query, uuid).Scan(&post.Id, &post.Uuid, &post.Title, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt, &post.CreatedAt)
	return
}

// Create a post
func (user *User) CreatePost(post *Post) (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO posts (uuid,title,body,user_id,user_id_str,user_name)
			  VALUES (?,?,?,?,?,?)`
	_, err = db.Exec(query, uuid.New().String(), post.Title, post.Body, user.Id, user.UserIdStr, user.Name)
	return
}

// Update the post
func (post *Post) UpdatePost() (err error) {
	db := NewDB()
	defer db.Close()
	query := `UPDATE posts
			  SET title=?,body=?
			  WHERE id=?`
	_, err = db.Exec(query, post.Title, post.Body, post.Id)
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
