package service

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		Directory:  "views",
		Extensions: []string{".gtpl"},
		Layout:     "layout",
	})

	mx := mux.NewRouter()
	initRoutes(mx, formatter)

	n := negroni.Classic()
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/login", LoginHandler(formatter)).Methods("GET", "POST")
	mx.HandleFunc("/upload", UploadHandler(formatter)).Methods("GET", "POST")
	mx.NotFoundHandler = NotFoundHandler(formatter)

	// 表示路由前缀为/views的请求都由该Handler处理
	mx.PathPrefix("/views").Handler(http.StripPrefix("/views", http.FileServer(http.Dir("/views"))))
}
