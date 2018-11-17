package service

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unrolled/render"
)

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

func NotFoundHandler(formatter *render.Render) http.HandlerFunc {
	// 此函数处理NotFound
	return func(w http.ResponseWriter, req *http.Request) {
		// 调用http.Error(http.ResponseWriter, error string, code int)
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	}
}
