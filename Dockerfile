FROM golang
MAINTAINER Karthik Gaekwad

# mgo is needed
RUN go get gopkg.in/mgo.v2

ADD . /go/src/github.com/karthequian/visitor
RUN go install github.com/karthequian/visitor
ENTRYPOINT /go/bin/visitor
EXPOSE 8080