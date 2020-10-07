package route

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
	"gopkg.in/go-playground/validator.v9"
)

// GET /post/new
// Show the form page to create the new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		tmp := parseTemplateFiles("layout", "post.new", "navbar.private")
		uuid := uuid.New()
		tmp.Execute(w, uuid)
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
	// uuid := r.PostFormValue("uuid")
	if err = post.Create(); err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/post/new", 302)
		return
	}
	url := fmt.Sprint("/post/show?id=", post.Uuid)
	http.Redirect(w, r, url, 302)
}

// GET /post/show
// Show the post
func ShowPost(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	post, err := data.PostByUuid(uuid)
	if err != nil {
		log.Fatal(err)
	}
	session, err := session(w, r)
	if session.UserId != post.UserId {
		tmp := parseTemplateFiles("layout", "navbar.private", "post.show.public")
		if err := tmp.Execute(w, post); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	} else {
		tmp := parseTemplateFiles("layout", "navbar.private", "post.show.private")
		if err := tmp.Execute(w, post); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	}
}

// GET /post/edit
// Show the post edit form
func EditPost(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	post, _ := data.PostByUuid(uuid)
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		tmp := parseTemplateFiles("layout", "navbar.private", "post.edit")
		if err := tmp.Execute(w, post); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	}
}

// POST /post/update
// Update the post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		user, err := sess.User()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		uuid := r.FormValue("uuid")
		post, err := data.PostByUuid(uuid)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		if user.Id != post.UserId {
			http.Redirect(w, r, "/", 302)
		}
		// attr := map[string]interface{}{
		// 	"Title": r.PostFormValue("title"),
		// 	"Body":  r.PostFormValue("body"),
		// }
		post.Title = r.PostFormValue("title")
		post.Body = r.PostFormValue("body")
		if err != nil {
			tmp := parseTemplateFiles("layout", "navbar.private", "post.edit")
			if err := tmp.Execute(w, post); err != nil {
				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			}
		} else {
			err := post.UpdatePost()
			if err != nil {
				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			}
			http.Redirect(w, r, "/", 302)
		}
	}
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
