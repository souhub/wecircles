package route

import (
	"fmt"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func CreateChat(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	url := fmt.Sprintf("/circle?id=%s", myUser.UserIdStr)
	user, err := data.UserByUserIdStr(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	circle, err := data.CirclebyOwnerID(user.UserIdStr)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, url, 302)
		return
	}
	chat := data.Chat{
		Body:             r.FormValue("body"),
		UserID:           myUser.Id,
		UserIdStr:        myUser.UserIdStr,
		UserImagePath:    myUser.ImagePath,
		CircleID:         circle.ID,
		CircleOwnerIDStr: circle.OwnerIDStr,
	}
	if err := chat.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, url, 302)
		return
	}
	url = fmt.Sprintf("/circle?id=%s", user.UserIdStr)
	http.Redirect(w, r, url, 302)
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	circleOwnerID := vals.Get("ownerid")
	_, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	chat, err := data.GetChatByID(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if err := chat.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	url := fmt.Sprintf("/circle?id=%s", circleOwnerID)
	http.Redirect(w, r, url, 302)
}
