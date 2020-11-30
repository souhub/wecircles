package data

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/souhub/wecircles/pkg/logging"
)

type Post struct {
	Id            int
	Uuid          string
	Title         string
	Body          string
	ThumbnailPath string
	UserId        int
	UserIdStr     string
	UserName      string
	UserImagePath string
	CreatedAt     string
}

func S3Upload(imagePath string) error {
	// credentialsの作成
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	cfg := aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		// Endpoint:    aws.String("http://127.0.0.1:9000"),
	}
	// sessionの作成
	sess, err := session.NewSession(&cfg)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// ファイルを開く
	targetFilePath := imagePath
	f, err := os.Open(targetFilePath)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	defer f.Close()

	bucketName := os.Getenv("WECIRCLES_S3_IMAGE_BUCKET")
	objectKey := targetFilePath

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   f,
	})
	return err
}

func S3Delete(imagePath string) error {
	// credentialsの作成
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	cfg := aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		// Endpoint:    aws.String("http://127.0.0.1:9000"),
	}
	// sessionの作成
	sess, err := session.NewSession(&cfg)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	bucketName := os.Getenv("WECIRCLES_S3_IMAGE_BUCKET")
	objectKey := imagePath

	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectKey)})

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	return err
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
		err = rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.ThumbnailPath, &post.UserId, &post.UserIdStr, &post.UserName, &post.UserImagePath, &post.CreatedAt)
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
	err = db.QueryRow(query, uuid).Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.ThumbnailPath, &post.UserId, &post.UserIdStr, &post.UserName, &post.UserImagePath, &post.CreatedAt)
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
	err = db.QueryRow(query, id).Scan(&post.Id, &post.Uuid, &post.Title, &post.Body, &post.ThumbnailPath, &post.UserId, &post.UserIdStr, &post.UserName, &post.UserImagePath, &post.CreatedAt)
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
	query := `INSERT INTO posts (uuid, title, body, thumbnail_path, user_id, user_id_str, user_name, user_image_path, created_at)
			  VALUES (?,?,?,?,?,?,?,?,?)`
	_, err = db.Exec(query, post.Uuid, post.Title, post.Body, post.ThumbnailPath, post.UserId, post.UserIdStr, post.UserName, post.UserImagePath, post.CreatedAt)
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
	// Make thumbnail dir.
	currentRootDir, err := os.Getwd()
	thumbnailImageDir := fmt.Sprintf("%s/web/img/user%d/posts/post%s", currentRootDir, post.UserId, post.Uuid)
	if _, err = os.Stat(thumbnailImageDir); err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		if err = os.MkdirAll(thumbnailImageDir, 0777); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
	}
	// Parse the form
	err = r.ParseForm()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Get the file sent form the form
	file, fileHeader, err := r.FormFile("image")
	// Get the uploaded file's name from the file.
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		uploadedFileName = post.ThumbnailPath
		return
	}
	uploadedFileName = fileHeader.Filename

	// Set the uploaded file's path
	imagePath := fmt.Sprintf("web/img/user%d/posts/post%s/%s", post.UserId, post.Uuid, uploadedFileName)

	// Save the uploaded file to "imagePath"
	saveImage, err := os.Create(imagePath)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Write the uploaded file to the file for saving on the server.
	_, err = io.Copy(saveImage, file)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	// Close the "saveImage" and "file"
	saveImage.Close()
	file.Close()

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

// Delete the thumbnail
func (post *Post) DeleteThembnail() (err error) {
	// currentRootDir, err := os.Getwd()
	thumbnail := fmt.Sprintf("web/img/user%d/posts/post%s/%s", post.UserId, post.Uuid, post.ThumbnailPath)
	// if _, err = os.Stat(thumbnail); err != nil {
	// 	logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	// 	return
	// }
	// err = os.Remove(thumbnail)
	if err = S3Delete(thumbnail); err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
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
