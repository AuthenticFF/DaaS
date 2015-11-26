FROM golang:latest

# Install webkit/gtk 
RUN apt-get update && apt-get install -y  --fix-missing libwebkit2gtk-3.0-dev xvfb

#configure xvfb
ENV DISPLAY :99
ADD xvfb-init /etc/init.d/xvfb
RUN chmod a+x /etc/init.d/xvfb
ADD ./init.sh /init.sh
RUN chmod a+x /init.sh

#build Go App
RUN go get github.com/Ramshackle-Jamathon/DaaS

ADD . /go/src/DaaS
WORKDIR /go/src/DaaS
RUN go get
RUN go install


ENV PORT=8080
EXPOSE 8080

#crush it
ENTRYPOINT /init.sh go run server.go