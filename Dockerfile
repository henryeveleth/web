FROM golang:latest
LABEL maintainer="Henry Eveleth <henryeveleth@gmail.com>"
RUN mkdir /app
COPY ./ $GOPATH/app
WORKDIR $GOPATH/app
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux
RUN go get github.com/sirupsen/logrus
RUN go build
EXPOSE 8080
CMD ["./web"]
