package route

import (
	"log"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
)

func index(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	posts := data.Posts()
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
