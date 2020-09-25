package route

import (
	"log"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

// GET /
// Show all of the posts
func Index(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	posts, err := data.Posts()
	if err != nil {
		logging.Warn("Failed to get all of the posts.")
	}
	if err != nil {
		tmp := parseTemplateFiles("layout", "index", "navbar.public")
		if err := tmp.Execute(w, posts); err != nil {
			logging.Warn("Failed to parse the templates.")
		}
	} else {
		tmp := parseTemplateFiles("layout", "index", "navbar.private")
		if err := tmp.Execute(w, posts); err != nil {
			logging.Warn("Failed to parse the templates.")
		}
	}
}
