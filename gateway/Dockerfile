FROM golang:latest

WORKDIR /go/src/app

COPY src/ .

RUN go build -o main .

CMD ["./main"]
