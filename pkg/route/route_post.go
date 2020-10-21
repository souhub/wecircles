package route

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
	"gopkg.in/go-playground/validator.v9"
)

// GET /post
func Posts(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	posts, err := user.PostsByUser()
	type Data struct {
		User  data.User
		Posts []data.Post
	}
	data := Data{
		User:  user,
		Posts: posts,
	}
	tmp := parseTemplateFiles("layout", "index", "navbar.private")
	tmp.Execute(w, data)
}

// GET /post/new
// Get the form page to create the new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	type Data struct {
		User data.User
		UUID uuid.UUID
	}
	data := Data{
		User: user,
		UUID: uuid.New(),
	}
	tmp := parseTemplateFiles("layout", "post.new", "navbar.private")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// POST /post/create
// Create the new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	err = r.ParseForm()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	user, err := sess.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	//createdAtを作成
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	date := now.Day()
	hour := now.Hour()
	createdAt := fmt.Sprintf("%v年%v月%v日%v時", year, month, date, hour)
	//postを作成
	post := data.Post{
		Uuid:      r.PostFormValue("uuid"),
		Title:     r.PostFormValue("title"),
		Body:      r.PostFormValue("body"),
		UserId:    user.Id,
		UserIdStr: user.UserIdStr,
		UserName:  user.Name,
		CreatedAt: createdAt,
	}
	validate := validator.New()  //validatorインスタンス生成
	err = validate.Struct(&post) //validator実行
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	if err = post.Create(); err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	url := fmt.Sprint("/post/show?id=", post.Uuid)
	http.Redirect(w, r, "/", 302)
	// この投稿にこのurlを割り当てるためのダミー（エラー発生するが問題ない）
	http.Redirect(w, r, url, 302)
}

// GET /post/show
// Get the post
func ShowPost(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	post, err := data.PostByUuid(id)
	if err != nil {
		log.Fatal(err)
	}
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	// ログイン判定
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		tmp := parseTemplateFiles("layout", "navbar.public", "post.show")
		if err := tmp.Execute(w, post); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		return
	}
	data := Data{
		User: user,
		Post: post,
	}
	// 投稿者判定はテンプレートで行う
	tmp := parseTemplateFiles("layout", "navbar.private", "post.show")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// GET /post/edit
// Get the post edit form
func EditPost(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	post, err := data.PostByUuid(uuid)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	data := Data{
		User: user,
		Post: post,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "post.edit")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// GET /post/edit/thumbnail
// Get the thumbnail edit form
func EditPostThumbnail(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	vals := r.URL.Query()
	uuid := vals.Get("id")
	post, err := data.PostByUuid(uuid)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "post.edit.thumbnail")
	if err := tmp.Execute(w, post); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

// POST /post/update
// Update the post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	vals := r.URL.Query()
	id := vals.Get("id")
	post, err := data.PostByUuid(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	post.Title = r.PostFormValue("title")
	post.Body = r.PostFormValue("body")
	if session.UserId != post.UserId {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/", 302)
		return
	}
	if err != nil {
		http.Redirect(w, r, "/post/edit", 302)
		return
	}
	if err := post.UpdatePost(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	http.Redirect(w, r, "/", 302)
}

// POST /post/update/thumbnail
// Update the thumbnail
func UpdatePostThumbnail(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	// 画像を multiple で送信しているため PaeseForm での解析不可
	err = r.ParseMultipartForm(0)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	uuid := r.FormValue("uuid")
	post, err := data.PostByUuid(uuid)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// userごとの画像保存フォルダにサムネイル保存フォルダ作成
	currentRootDir, err := os.Getwd()
	thumbnailImageDir := fmt.Sprintf("%s/web/img/user%d/posts/post%d", currentRootDir, post.UserId, post.Id)
	_, err = os.Stat(thumbnailImageDir)
	if err != nil {
		err = os.MkdirAll(thumbnailImageDir, 0777)
	}
	// Delete the current user's image
	err = post.DeleteThembnail()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// Uplaod the new user's image.
	thumbnailPath, err := post.UploadThumbnail(r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	post.ThumbnailPath = thumbnailPath
	// Update the user's image path from old one to the new one in the DB.
	if err = post.UpdateThumbnail(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// url := fmt.Sprintf("/post/edit?id=%s", uuid)
	http.Redirect(w, r, "/", 302)
}

// DELETE /post/delete
// Delete the post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	vals := r.URL.Query()
	uuid := vals.Get("id")
	user, err := sess.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	post, err := data.PostByUuid(uuid)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	if user.Id != post.UserId {
		http.Redirect(w, r, "/", 302)
	}
	err = post.Delete()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	http.Redirect(w, r, "/", 302)

}
