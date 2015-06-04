FROM golang

ADD . /go/src/github.com/tobyhughes/shortie

RUN go install github.com/tobyhughes/shortie

ENTRYPOINT /go/bin/shortie
