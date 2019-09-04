FROM golang:latest
LABEL maintainer="Henry Eveleth <henryeveleth@gmail.com>"
RUN mkdir -p $GOPATH/src/github.com/henryeveleth/web/
ADD . $GOPATH/src/github.com/henryeveleth/web/
WORKDIR $GOPATH/src/github.com/henryeveleth/web/
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux
RUN go get github.com/sirupsen/logrus
RUN go build .
EXPOSE 8080
CMD ["./web"]
