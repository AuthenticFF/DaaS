FROM golang:latest


ADD ./init.sh /init.sh
RUN chmod a+x /init.sh

#build Go App
RUN go get github.com/AuthenticFF

ADD . /go/src/github.com/AuthenticFF/DaaS
WORKDIR /go/src/github.com/AuthenticFF/DaaS
RUN go get
RUN go install


ENV PORT=9091
EXPOSE 9091

#crush it
ENTRYPOINT /init.sh go run server.go

#excellent command
#docker rmi -f $(docker images | grep "^<none>" | awk "{print $3}")