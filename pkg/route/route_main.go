package route

import (
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

// GET
// Show all of the posts
func Index(w http.ResponseWriter, r *http.Request) {
	posts, err := data.Posts()
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		Posts: posts,
	}
	session, err := session(w, r)
	if err != nil {
		data := Data{
			Posts: posts,
		}
		tmp := parseTemplateFiles("layout", "index", "navbar.public", "posts")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		return
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	data = Data{
		MyUser: myUser,
		Posts:  posts,
	}
	tmp := parseTemplateFiles("layout", "index", "navbar.private", "posts")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}
