package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

// GET /circle/memberships
// Show the circles memberships
func MembershipsCircles(w http.ResponseWriter, r *http.Request) {
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
	memberships, err := myUser.MembershipsByUserID()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	var circles []data.Circle
	for _, membership := range memberships {
		circle, err := membership.Circle()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		circles = append(circles, circle)
	}
	data := Data{
		MyUser:          myUser,
		Circles:         circles,
		ImagePathPrefix: os.Getenv("IMAGE_PATH"),
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circles")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

// POST /circle/membership/create
// Create the membership.
func MembershipsCircleCreate(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	circle, err := data.CirclebyOwnerID(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	membership := data.Membership{
		UserID:   user.Id,
		CircleID: circle.ID,
	}
	if err := membership.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	url := fmt.Sprintf("/circle?id=%s", circle.OwnerIDStr)
	http.Redirect(w, r, url, 302)
}

// DELETE /circle/membership/delete
// Delete the membership.
func DeleteMembership(w http.ResponseWriter, r *http.Request) {
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
	circle, err := data.CirclebyOwnerID(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	membership := data.Membership{
		UserID:   myUser.Id,
		CircleID: circle.ID,
	}
	if err := membership.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	url := fmt.Sprintf("/circle?id=%s", circle.OwnerIDStr)
	http.Redirect(w, r, url, 302)
}

// DELETE /circle/membership/delete
// Delete the membership.
func DeleteMembershipByOwner(w http.ResponseWriter, r *http.Request) {
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
	circle, err := data.CirclebyOwnerID(myUser.UserIdStr)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	deletedUser, err := data.UserByUserIdStr(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	membership := data.Membership{
		UserID:   deletedUser.Id,
		CircleID: circle.ID,
	}
	if err := membership.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	http.Redirect(w, r, "/circle/manage/members", 302)
}
