package main

import (
	"net/http"

	"github.com/souhub/wecircles/pkg/route"
)

func main() {
	http.HandleFunc("/", route.Index)
	http.HandleFunc("/show", route.Show)
	http.HandleFunc("/calc", route.Calc)
	http.HandleFunc("/error", route.Error)
	http.ListenAndServe(":80", nil)
}
