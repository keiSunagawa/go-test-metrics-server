FROM golang:latest

WORKDIR /go
ADD . /go

RUN go get github.com/keiSunagawa/go-test-metrics-server

CMD ["go-test-metrics-server"]
