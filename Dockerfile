FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

EXPOSE 80

CMD ["go", "run", "main.go"]