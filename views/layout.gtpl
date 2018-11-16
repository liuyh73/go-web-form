{{ define "layout" }}
<html>
  <head>
    <title>go-web-form</title>
    <link rel="icon" href="/static/images/icon.png" />
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