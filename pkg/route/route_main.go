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
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	_, err = session(w, r)
	if err != nil {
		tmp := parseTemplateFiles("layout", "index", "navbar.public")
		if err := tmp.Execute(w, posts); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	} else {
		tmp := parseTemplateFiles("layout", "index", "navbar.private")
		if err := tmp.Execute(w, posts); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	}
}
