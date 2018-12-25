FROM golang:latest
MAINTAINER liuyh73  "15989067460@163.com"
# 该指令用于配置工作目录，其参数应该使用绝对目录。
WORKDIR $GOPATH/src/github.com/liuyh73/go-web-form
# ADD不但支持将本地文件复制到容器中，还支持本地提取文件和远程url下载
# ADD <src> <dst>
ADD . $GOPATH/src/github.com/liuyh73/go-web-form
RUN go get github.com/gorilla/mux
RUN go get github.com/codegangsta/negroni
RUN go get github.com/unrolled/render
RUN go get github.com/spf13/pflag
RUN go build .
# 该指令指示容器讲监听链接的端口，类似于，将容器中的某一个端口暴露出去，从而在外部访问绑定该端口。
EXPOSE 8080
# ENTRYPOINT允许你配置作为可执行文件运行的容器
ENTRYPOINT ["./go-web-form"]
