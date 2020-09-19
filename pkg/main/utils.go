package main

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"
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
		files = append(files, fmt.Sprintf("view/templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}
