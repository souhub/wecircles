package route

import (
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
	"gopkg.in/go-playground/validator.v9"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tmp := parseTemplateFiles("login.layout", "navbar.public", "login")
	err := tmp.Execute(w, nil)
	if err != nil {
		logging.Warn("Failet to open Login page.")
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	tmp := parseTemplateFiles("login.layout", "navbar.public", "signup")
	err := tmp.Execute(w, nil)
	if err != nil {
		logging.Warn("Failet to open Signup page.")
	}
}

func SignupAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logging.Warn("Failed to signup, because of parsing forms.")
		http.Redirect(w, r, "/signup", 302)
		return
	}
	user := data.User{
		Name:      r.PostFormValue("name"),
		UserIdStr: r.PostFormValue("user_id_str"),
		Email:     r.PostFormValue("email"),
		Password:  data.Encrypt(r.PostFormValue("password")),
	}
	validate := validator.New() //validatorインスタンス生成
	//validator実行
	if err := validate.Struct(&user); err != nil {
		logging.Warn("Failed to create user and redirect to the Signup page, because of validation.")
		http.Redirect(w, r, "/signup", 302)
		return
	}
	//DBに登録
	if err := user.Create(); err != nil {
		logging.Warn("Failed to signup and redirect to the Signup page.")
		http.Redirect(w, r, "/signup", 302)
		return
	}
	http.Redirect(w, r, "/login", 302)
	//そのまま認証終わらせてマイページに飛ばす
	// session, err := user.CreateSession()
	// if err != nil {
	// 	logging.Warn("Failed to create the session and redirect to Signup page")
	// 	http.Redirect(w, r, "/signup", 302)
	// 	return
	// }
	// cookie := http.Cookie{
	// 	Name:     "_cookie",
	// 	Value:    session.Uuid,
	// 	HttpOnly: true,
	// }
	// http.SetCookie(w, &cookie)
	// url := fmt.Sprint("/user/show?id=", session.UserIdStr)
	// http.Redirect(w, r, url, 302)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	//フォームの入力内容解析
	if err := r.ParseForm(); err != nil {
		logging.Warn("Failed to parse the form.")
		http.Redirect(w, r, "/login", 302)
		return
	}
	//フォームに入力されたEmailを持つユーザーをデータベースから特定

	user, err := data.UserByUserIdStr(r.FormValue("useridstr"))
	if err != nil {
		logging.Warn("Failed to find user by the User IDddddddddd.")
		http.Redirect(w, r, "/login", 302)
		return
	}
	//その特定したユーザーとフォームに入力されたパスワードのチェック
	if user.Password == data.Encrypt(r.FormValue("password")) {
		//認証されたらユーザーのセッションを作成しクッキーを持たせ、マイページにリダイレクトさせる
		//SHA-1でuuid作ってクッキーを作ってブラウザに渡す。
		session, err := user.CreateSession()
		if err != nil {
			log.Warn("Failed to create the new session")
			http.Redirect(w, r, "/login", 302)
			return
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		//認証されなかったらログインファームにリダイレクトさせる
		http.Redirect(w, r, "/login", 302)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.Delete()
	}
	http.Redirect(w, r, "/login", 302)
}
