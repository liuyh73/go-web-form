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