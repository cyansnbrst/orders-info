FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server /app/cmd/api

EXPOSE 8080

CMD ["/app/server"]
