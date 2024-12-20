FROM golang:1.22.7 AS builder

WORKDIR /client

COPY ../ ./

CMD ["go", "run", "client.go"]
