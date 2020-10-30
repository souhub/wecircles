package route

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
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

func checkMembership(user data.User, circle data.Circle) (ok bool, err error) {
	memberships, err := user.Memberships()
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	for _, membership := range memberships {
		ok, err = membership.Check(circle)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		if !ok {
			err = errors.New("Invalid the membership")
			return
		}
		return
	}
	return
}

type Data struct {
	MyUser          data.User
	User            data.User
	Users           []data.User
	Post            data.Post
	Posts           []data.Post
	Circle          data.Circle
	Circles         []data.Circle
	Membership      data.Membership
	Memberships     []data.Membership
	Chat            data.Chat
	Chats           []data.Chat
	MembershipValid bool
}
