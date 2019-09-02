
FROM golang:1.12.9

LABEL maintainer="Bubunyo Nyavor <kiddbubu@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/goose ./cmd/goose/main.go

RUN go build -o ./bin/application ./*.go

EXPOSE 8080

CMD ["./bin/application"]
