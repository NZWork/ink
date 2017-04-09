FROM golang:1.7

ADD . /go/src/ink

WORKDIR /go/src/ink
# RUN godep go install ink
COPY vendor/ /go/src/

RUN go install ink

EXPOSE 8080

ENTRYPOINT /go/bin/ink
