package route

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

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

// Create the user image
func upload(w http.ResponseWriter, r *http.Request) {
	//メソッドをPOSTのみ許可
	if r.Method != "POST" {
		log.Fatalln("Panic")
	}

	//formから送信されたファイルを解析
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	//アップロードされたファイル名を取得
	uploadedFileName := fileHeader.Filename
	//アップロードされたファイルを置くパスを設定
	imagePath := "web/img/user/" + uploadedFileName

	//imagePathにアップロードされたファイルを保存
	saveImage, err := os.Create(imagePath)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	//保存用ファイルにアップロードされたファイルを書き込む
	_, err = io.Copy(saveImage, file)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	//saveImageとfileを最後に閉じる
	defer saveImage.Close()
	defer file.Close()

	//もう1周
	http.Redirect(w, r, "/mypage", 302)
}
