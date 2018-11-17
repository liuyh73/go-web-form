## 服务搭建
### server.go
```go
// 创建服务函数，返回negroni.Negroni指针
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
```
**initRouters()**函数如下：
```go
func initRoutes(mx *mux.Router, formatter *render.Render) {
    // 调用os.Getwd()获取目录, 用于后面静态资源定位
	path, _ := os.Getwd()
    // 注册路由，处理Methods：GET和POST
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
```
### router.go
**LoginHandler**
```go
// 定义路由处理函数
func LoginHandler(formatter *render.Render) http.HandlerFunc {
	// 返回http.HandlerFunc,处理GET和POST请求
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("method:", req.Method)
		// formatter为一个渲染模板的render实例
		// formatter.HTML(http.ResponseWriter, http.StatusCode, HTML模板, 模板绑定的值)
		if req.Method == "GET" {
			formatter.HTML(w, http.StatusOK, "layout", true)
		} else {
			// req.ParseForm()获取表单提交的值
			req.ParseForm()
			// 自定义模板，可以使用ParseFiles利用模板文件获取template.Template对象
			// {{define "username"}} …… {{end}} 给模板命名
			t, _ := template.New("login").Parse(`{{define "username"}}Hello, {{.}}!{{end}}`)
			// t.ExecuteTemplate(http.ResponseWriter, 模板名称, 模板对象的值)
			log.Println(t.ExecuteTemplate(w, "username", req.Form.Get("username")))
		}
	}
}
```
**UploadHandler**
```go
func UploadHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("method: ", req.Method)
		if req.Method == "GET" {
			// Get方法渲染upload模板
			formatter.HTML(w, http.StatusOK, "layout", false)
		} else {
			// 上传文件是需要调用req.ParseMultipartForm, 参数为最大占用存储空间,将request body转化为multipart/form-data,
			req.ParseMultipartForm(32 << 20)
			// 获取文件, req.FormFile("")参数为input表单的name属性
			file, handler, err := req.FormFile("uploadfile")
			defer file.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Fprintf(w, "%v", handler.Header)
			// 将上传文件拷贝到本地
			f, err := os.OpenFile("./file/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			io.Copy(f, file)
		}
	}
}
```
**NotFoundHandler**
```go
func NotFoundHandler(formatter *render.Render) http.HandlerFunc {
	// 此函数处理NotFound
	return func(w http.ResponseWriter, req *http.Request) {
		// 调用http.Error(http.ResponseWriter, error string, code int)
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	}
}
```
