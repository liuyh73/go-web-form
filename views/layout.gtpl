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