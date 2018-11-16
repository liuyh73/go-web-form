{{ define "upload" }}
  <form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="submit" value="upload" />
  </form>
{{ end }}