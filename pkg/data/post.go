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

func Posts() (posts []Post, err error) {
	defer db.Close()
	cmd := "SELECT * FROM posts"
	rows, err := db.Query(cmd)
	if err != nil {
		logging.Warn("Coudn't find posts.")
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt, &post.CreatedAt)
		if err != nil {
			logging.Warn("Couldn't find a post.")
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func PostByUuid(uuid string) (post Post, err error) {
	post = Post{}
	defer db.Close()
	cmd := "SELECT * FROM users WHERE uuid=$1"
	err = db.QueryRow(cmd, uuid).Scan(&post.Id, &post.Uuid, &post.Title, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt, &post.CreatedAt)
	return
}

func (user *User) CreatePost(post *Post) (err error) {
	defer db.Close()
	cmd := "INSERT INTO posts (uuid,title,body,user_id,user_id_str,user_name) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err = db.Exec(cmd, uuid.New().String(), post.Title, post.Body, user.Id, user.UserIdStr, user.Name)
	return
}

func (post *Post) UpdatePost() (err error) {
	defer db.Close()
	cmd := "UPDATE posts SET title=$2,body=$3 WHERE id=$1"
	_, err = db.Exec(cmd, post.Id, post.Title, post.Body)
	return
}

func (post *Post) DeletePost() (err error) {
	defer db.Close()
	cmd := "DELETE from posts WHERE id=$1"
	_, err = db.Exec(cmd, post.Id)
	return
}

func (user *User) PostsByUser() (posts []Post, err error) {
	defer db.Close()
	cmd := "SELECT * FROM posts WHERE user_id=$1"
	rows, err := db.Query(cmd, user.Id)
	if err != nil {
		logging.Warn("Coudn't find posts.")
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.UserId, &post.UserIdStr, &post.UserName, &post.CreatedAt, &post.CreatedAt)
		if err != nil {
			logging.Warn("Couldn't find a post.")
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}
