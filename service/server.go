package service

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func NewServer() *negroni.Negroni {
	// 返回一个Render实例的指针，Render是一个包，提供轻松呈现JSON，XML，文本，二进制数据和HTML模板的功能
	// Directory : Specify what path to load the templates from.
	// Layout : Specify a layout template. Layouts can call {{ yield }} to render the current template or {{ partial "css" }} to render a partial from the current template.
	// Extensions: Specify extensions to load for templates.
	formatter := render.New(render.Options{
		Directory:  "views",
		Extensions: []string{".gtpl"},
		Layout:     "layout",
	})
	// 设置router
	mx := mux.NewRouter()
	initRoutes(mx, formatter)

	// negroni.Classic() 返回带有默认中间件的Negroni实例指针:
	n := negroni.Classic()
	// 让 negroni 使用该 Router
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	// 注册路由，处理Methods：GET和POST
	path, _ := os.Getwd()
	mx.HandleFunc("/login", LoginHandler(formatter)).Methods("GET", "POST")
	mx.HandleFunc("/upload", UploadHandler(formatter)).Methods("GET", "POST")
	mx.NotFoundHandler = NotFoundHandler(formatter)

	// 表示路由前缀为"/views"的请求都由该Handler处理
	// mx.PathPrefix("")匹配前缀，返回*mux.Route, 链式调用Handler(http.Handler)
	// http.StripPrefix("", http.Handler)去除前缀, 并将请求定向到http.Handler
	// http.FileServer(http.FileSystem) 返回http.Handler
	// http.Dir("")参数应该为绝对路径
	mx.PathPrefix("/views").Handler(http.StripPrefix("/views", http.FileServer(http.Dir(path+"/views"))))
	mx.PathPrefix("/static/images").Handler(http.StripPrefix("/static/images", http.FileServer(http.Dir(path+"/static/images"))))
}
