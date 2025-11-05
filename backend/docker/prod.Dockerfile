FROM golang:1.24.6-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 80

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /tmp/main .

CMD ["/tmp/main"]