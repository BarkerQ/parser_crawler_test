FROM golang:1.17

WORKDIR /go/src/parserProject
COPY . .
COPY .env /go/src/parserProject/main
WORKDIR /go/src/parserProject/main

RUN go build main.go
ENTRYPOINT ["./main"]
