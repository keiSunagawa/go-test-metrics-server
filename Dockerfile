FROM golang:latest

WORKDIR /go
ADD . /go

RUN go get github.com/keiSunagawa/go-test-metrics-server
#RUN go get -u github.com/prometheus/client_golang/master
CMD ["go-test-metrics-server"]
