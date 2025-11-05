FROM golang:1.24.6-alpine AS base

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/air-verse/air@v1.62.0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN chmod +x ./entrypoint.sh

EXPOSE 80

FROM base AS dev

ENV GO_PROD = 0

COPY --from=base /go/bin/goose /usr/local/bin/goose
COPY --from=base /go/bin/air /usr/local/bin/air
COPY --from=base /app/entrypoint.sh /usr/local/bin/entrypoint

CMD [ "entrypoint" ]

FROM base AS prod

ENV GO_PROD = 1

COPY --from=base /go/bin/goose /usr/local/bin/goose

COPY --from=base /app/entrypoint.sh /usr/local/bin/entrypoint

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /tmp/main .

CMD [ "entrypoint" ]