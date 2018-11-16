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

func LoginHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("method:", req.Method)
		if req.Method == "GET" {
			formatter.HTML(w, http.StatusOK, "login", true)
		} else {
			req.ParseForm()
			t, _ := template.New("login").Parse(`{{define "username"}}Hello, {{.}}!{{end}}`)
			log.Println(t.ExecuteTemplate(w, "username", req.Form.Get("username")))
		}
	}
}

func UploadHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("method: ", req.Method)
		if req.Method == "GET" {
			formatter.HTML(w, http.StatusOK, "upload", false)
		} else {
			req.ParseMultipartForm(32 << 20)
			file, handler, err := req.FormFile("uploadfile")
			defer file.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Fprintf(w, "%v", handler.Header)

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
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	}
}
