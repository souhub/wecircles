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

func NewPost(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// tmp := template.Must(template.ParseFiles("templates/layout.html", "templates/post.new.html", "templates/navbar.private.html"))
		tmp := parseTemplateFiles("layout", "post.new", "navbar.private")
		uuid := uuid.New()
		tmp.Execute(w, uuid)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			logging.Warn("NoForm")
		}
		user, err := sess.User()
		if err != nil {
			logging.Warn("NoUser")
		}
		//createdAtを作成
		now := time.Now()
		year := now.Year()
		month := int(now.Month())
		date := now.Day()
		hour := now.Hour()
		createdAt := fmt.Sprintf("%v年%v月%v日 %v時", year, month, date, hour)
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
			http.Redirect(w, r, "/post/new", 302)
			return
		}
		uuid := r.PostFormValue("uuid")
		if err = user.CreatePost(&post); err != nil {
			log.Fatal(err)
		}
		url := fmt.Sprint("/post/show?id=", uuid)
		http.Redirect(w, r, url, 302)

	}
}

func ShowPost(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	post, _ := data.PostByUuid(uuid)
	_, err := session(w, r)
	if err != nil {
		tmp := parseTemplateFiles("layout", "navbar.public", "post.show.public")
		if err := tmp.Execute(w, post); err != nil {
			log.Fatal(err)
		}
	} else {
		tmp := parseTemplateFiles("layout", "navbar.private", "post.show.private")
		if err := tmp.Execute(w, post); err != nil {
			log.Fatal(err)
		}
	}
}

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
			log.Fatal(err)
		}
	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			logging.Warn("NoUuid")
		}
		user, err := sess.User()
		if err != nil {
			logging.Warn("NoUser")
		}
		uuid := r.FormValue("uuid")
		post, err := data.PostByUuid(uuid)
		if err != nil {
			logging.Warn("NoUuid")
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
				logging.Warn("Failed to exec templates.")
			}
		} else {
			err := post.UpdatePost()
			if err != nil {
				logging.Warn("Failed to update your post.")
			}
			http.Redirect(w, r, "/", 302)
		}
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	vals := r.URL.Query()
	uuid := vals.Get("id")
	user, err := sess.User()
	if err != nil {
		logging.Warn("Failed to find your user account from the session.")
	}
	post, err := data.PostByUuid(uuid)
	if err != nil {
		logging.Warn("Failed to find your post from the uuid of the post.")
	}
	if user.Id != post.UserId {
		http.Redirect(w, r, "/", 302)
	}
	err = post.DeletePost()
	if err != nil {
		logging.Warn("Failed to delete your post.")
	}
	http.Redirect(w, r, "/", 302)
}
