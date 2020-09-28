package data

import (
	"log"

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
		log.Fatal(err)
		logging.Warn("Failed to find posts.")
		return
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
			logging.Warn("Failed to find a post.")
			return
		}
		posts = append(posts, post)
	}
	defer rows.Close()
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
		log.Fatal(err)
		logging.Warn("Failed to find posts.")
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
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
			  FROM posts
			  WHERE uuid=?`
	err = db.QueryRow(query, uuid).Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt)
	return
}

// Create a post
func (post *Post) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO posts (uuid, title, body, user_id, user_id_str, user_name, created_at)
			  VALUES (?,?,?,?,?,?,?)`
	_, err = db.Exec(query, post.Uuid, post.Title, post.Body, post.UserId, post.UserIdStr, post.UserName, post.CreatedAt)
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
