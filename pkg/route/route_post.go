package route

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

// GET /post
func PostsManage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	myUser, err := session.User()
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	posts, err := myUser.PostsByUser()
	data := Data{
		MyUser:          myUser,
		Posts:           posts,
		ImagePathPrefix: os.Getenv("IMAGE_PATH"),
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "posts.manage")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

// GET /post/new
// Get the form page to create the new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	myUser, err := session.User()
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	type Data struct {
		MyUser          data.User
		UUID            uuid.UUID
		ImagePathPrefix string
	}
	data := Data{
		MyUser:          myUser,
		UUID:            uuid.New(),
		ImagePathPrefix: os.Getenv("IMAGE_PATH"),
	}
	tmp := parseTemplateFiles("post.new")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// POST /post/create
// Create the new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if err = r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	date := now.Day()
	hour := now.Hour()
	createdAt := fmt.Sprintf("%v年%v月%v日%v時", year, month, date, hour)
	if err := r.ParseMultipartForm(0); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	post := data.Post{
		Uuid:          uuid.New().String(),
		Title:         r.PostFormValue("title"),
		Body:          r.PostFormValue("body"),
		UserId:        user.Id,
		UserIdStr:     user.UserIdStr,
		UserName:      user.Name,
		UserImagePath: user.ImagePath,
		CreatedAt:     createdAt,
	}
	thumbnailPath, err := post.UploadThumbnail(r)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		thumbnailPath = "default_thumbnail.jpg"
	}
	post.ThumbnailPath = thumbnailPath
	if err = post.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	http.Redirect(w, r, "/", 302)
}

// GET /post/show
// Get the post
func ShowPost(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	post, err := data.PostByUuid(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	myUser, err := session.User()
	// ログイン判定
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		tmp := parseTemplateFiles("layout", "navbar.public", "post")
		if err := tmp.Execute(w, post); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		return
	}
	data := Data{
		MyUser: myUser,
		Post:   post,
	}
	// Contributor determination is done by template.
	tmp := parseTemplateFiles("layout", "navbar.private", "post")
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	data := Data{
		MyUser: myUser,
		Post:   post,
	}
	tmp := parseTemplateFiles("post.edit")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
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
	thumbnailPath, err := post.UploadThumbnail(r)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// If the thumbnails have been updated, delete the old ones from the folder.
	if post.ThumbnailPath != thumbnailPath {
		if err := post.DeleteThembnail(); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	}
	post.ThumbnailPath = thumbnailPath
	if session.UserId != post.UserId {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/post/edit", 302)
		return
	}
	post.Title = r.PostFormValue("title")
	post.Body = r.PostFormValue("body")
	if err := post.UpdatePost(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	http.Redirect(w, r, "/", 302)
}

// DELETE /post/delete
// Delete the post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	vals := r.URL.Query()
	uuid := vals.Get("id")
	post, err := data.PostByUuid(uuid)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	if myUser.Id != post.UserId {
		http.Redirect(w, r, "/", 302)
	}
	if err = post.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	currentRootDir, err := os.Getwd()
	postPath := fmt.Sprintf("%s/web/img/user%d/posts/post%s", currentRootDir, myUser.Id, post.Uuid)
	if _, err = os.Stat(postPath); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	if err := os.RemoveAll(postPath); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	http.Redirect(w, r, "/", 302)
}
