FROM golang

RUN go get github.com/go-sql-driver/mysql

RUN go get github.com/jessevdk/go-flags

ADD . /go/src/github.com/tobyhughes/shortie

RUN go install github.com/tobyhughes/shortie

ENTRYPOINT /go/bin/shortie
