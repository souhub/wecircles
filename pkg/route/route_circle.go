package route

import (
	"net/http"

	"github.com/souhub/wecircles/pkg/logging"
)

func NewCircle(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	tmp := parseTemplateFiles("layout", "circle.new", "navbar.private")
	tmp.Execute(w, nil)
}
