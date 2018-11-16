package main

import (
    "fmt"
    "html/template"
    "io"
    "log"
    "net/http"
    "os"
)

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        log.Println(t.Execute(w, nil))
    } else {
        r.ParseForm()
        t, _ := template.New("login").Parse(`{{define "username"}}Hello, {{.}}!{{end}}`)
        log.Println(t.ExecuteTemplate(w, "username", r.Form.Get("username")))
    }
}

func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method: ", r.Method)
    if r.Method == "GET" {
        t, _ := template.ParseFiles("upload.gtpl")
        t.Execute(w, nil)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
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

func main() {
    http.HandleFunc("/login", login)
    http.HandleFunc("/upload", upload)
    http.ListenAndServe(":9090", nil)
}
