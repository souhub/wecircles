package route

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
)

//"_cookie"のValueとUuidと同じUuidを持つSessionを取得
func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid Session")
		}
	}
	return
}

func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("web/templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}
