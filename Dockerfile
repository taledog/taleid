FROM golang:1.12-alpine
MAINTAINER luwei <taledog@gmail.com>

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:$GOROOT/bin
ENV GO111MODULE on

RUN mkdir -p /taleid /taleid/etc /taleid/logs
WORKDIR /taleid

ADD . /go/src/github.com/taledog/taleid

# 编译
RUN cd /go/src/github.com/taledog/taleid && \
  go build -v -o /taleid/idHttp  cmd/idHttp/main.go && \
  go build -v -o /taleid/idRedis cmd/idRedis/main.go && \
# 配置
  cp etc/*.toml /taleid/etc

CMD ["/taleid/idRedis"]
