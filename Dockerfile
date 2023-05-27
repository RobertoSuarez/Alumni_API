FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV URLDB="http://localhost/database-dockerfile"

EXPOSE 3000

RUN go build -o ./app main.go

CMD ["./app"]