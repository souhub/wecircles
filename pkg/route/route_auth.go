package route

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
	"gopkg.in/go-playground/validator.v9"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tmp := parseTemplateFiles("login.layout", "navbar.public", "login")
	tmp.Execute(w, nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	tmp := parseTemplateFiles("login.layout", "navbar.public", "signup")
	tmp.Execute(w, nil)
}

func SignupAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logging.Warn("Failed to signup")
	}
	user := data.User{
		Name:      r.PostFormValue("name"),
		UserIdStr: r.PostFormValue("user_id_str"),
		Email:     r.PostFormValue("email"),
		Password:  data.Encrypt(r.PostFormValue("password")),
	}
	validate := validator.New()  //validatorインスタンス生成
	err = validate.Struct(&user) //validator実行
	if err != nil {
		http.Redirect(w, r, "/signup", 302)
	}
	//DBに登録
	db := data.NewDB()
	defer db.Close()
	if err := user.Create(); err != nil {
		http.Redirect(w, r, "/signup", 302)
	}
	//そのまま認証終わらせてマイページに飛ばす
	session, err := user.CreateSession()
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.Uuid,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := fmt.Sprint("/user/show?id=", session.UserIdStr)
	http.Redirect(w, r, url, 302)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	//フォームの入力内容解析
	err := r.ParseForm()
	if err != nil {
		logging.Warn("Failed to parse the form.")
	}
	//フォームに入力されたEmailを持つユーザーをデータベースから特定
	email := r.FormValue("email")
	user, err := data.UserByEmail(email)
	if err != nil {
		logging.Warn("Failed to find user by the email address.")
	}
	//その特定したユーザーとフォームに入力されたパスワードのチェック
	password := r.FormValue("password")
	if user.Password == data.Encrypt(password) {
		//認証されたらユーザーのセッションを作成しクッキーを持たせ、マイページにリダイレクトさせる
		//SHA-1でuuid作ってクッキーを作ってブラウザに渡す。
		session, err := user.CreateSession()
		if err != nil {
			log.Warn("Failed to create the new session")
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
