package route

import (
	"net/http"
	"text/template"

	"github.com/souhub/wecircles/pkg/config"
	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/math"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("web/templates/index.html"))
	tmp.Execute(w, nil)
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := config.NewDB()
	var users []data.User
	db.Table("users").Find(&users)
	tmp := template.Must(template.ParseFiles("web/templates/show.html"))
	tmp.Execute(w, users)
}

func Calc(w http.ResponseWriter, r *http.Request) {
	data := math.Avg(10, 20)
	tmp := template.Must(template.ParseFiles("web/templates/calc.html"))
	tmp.Execute(w, data)
}
