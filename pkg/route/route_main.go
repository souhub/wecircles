package route

import (
	"log"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		logging.Warn("Coudn't find your session.")
	}
	posts, err := data.Posts()
	if err != nil {
		tmp := parseTemplateFiles("layout", "index", "navbar.public")
		if err := tmp.Execute(w, posts); err != nil {
			log.Fatal(err)
		}
	} else {
		tmp := parseTemplateFiles("layout", "index", "navbar.private")
		if err := tmp.Execute(w, posts); err != nil {
			log.Fatal(err)
		}
	}
}
