# go-web-form
## 目录结构
```bash
E:.
├─file(存储上传到的文件)
├─service(后台服务设计)
│   ├─router.go #路由
│   └─server.go #启动后台服务
├─static(静态资源)
│  └─images
│       └─favicon.ico
└─views(界面)
    ├─layout.gtpl #layout
    ├─login.gtpl  #login界面
    └─upload.gtpl #upload界面
```
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
## 界面设计
**layout.gtpl**
login.gtpl和upload.gtpl共用该文件，其中当该模板传入的参数为true时，{{template "login"}}会将login.gtpl include到此处；当该模板传入的参数为false时，{{template "upload"}}会将upload.gtpl include到此处；
```html
{{ define "layout" }}
<html>
  <head>
    <title>go-web-form</title>
    <link rel="icon" href="/static/images/favicon.ico" />
  </head>
  <body>
    {{ if . }}
    {{ template "login" }}
    {{ else }}
    {{ template "upload" }}
    {{ end }}
  </body>
</html>
{{ end }}
```
**login.gtpl**
```html
{{define "login"}}
  <form action="/login" method="post">
    用户名：<input type="text" name="username" />
    密码：<input type="password" name="password" />
    <input type="submit" value="登陆" />
  </form>
{{end}}
```
**upload.gtpl**
```html
{{ define "upload" }}
  <form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="submit" value="upload" />
  </form>
{{ end }}
```
## curl测试
在Linux中curl是一个利用URL规则在命令行下工作的文件传输工具，可以说是一款很强大的http命令行工具。它支持文件的上传和下载，是综合传输工具，但按传统，习惯称url为下载工具。
常见参数介绍：
```
-A/--user-agent <string>          设置用户代理发送给服务器
-b/--cookie <name=string/file>    cookie字符串或文件读取位置
-c/--cookie-jar <file>            操作结束后把cookie写入到这个文件中
-C/--continue-at <offset>         断点续转
-D/--dump-header <file>           把header信息写入到该文件中
-e/--referer                      来源网址
-f/--fail                         连接失败时不显示http错误
-o/--output                       把输出写到该文件中
-O/--remote-name                  把输出写到该文件中，保留远程文件的文件名
-r/--range <range>                检索来自HTTP/1.1或FTP服务器字节范围
-s/--silent                       静音模式。不输出任何东西
-T/--upload-file <file>           上传文件
-u/--user <user[:password]>       设置服务器的用户和密码
-w/--write-out [format]           什么输出完成后
-x/--proxy <host[:port]>          在给定的端口上使用HTTP代理
-#/--progress-bar                 进度条显示当前的传送状态
```
-   `curl -v 127.0.0.1:8080/login` Method：GET
![](./screenshots/1.png)
-   `curl -v 127.0.0.1:8080/upload` Method：GET
![](./screenshots/2.png)
-   `curl -v 127.0.0.1:8080/login -X POST -d "username=liuyh73&&password=acwab"` Method：POST
![](./screenshots/3.png)
-   `curl -v 127.0.0.1:8080/upload -F "uploadfile=@E:/mygo/src/github.com/liuyh73/go-web-form/static/images/favicon.ico"` Method：POST
![](./screenshots/4.png)
-   `curl -v 127.0.0.1:8080/static/images/favicon.ico` Method: GET
![](./screenshots/5.png)
## ab测试
