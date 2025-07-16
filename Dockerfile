FROM golang:1.24.4-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/fin-agg ./cmd/main
RUN CGO_ENABLED=0 go build -o /app/migrate ./cmd/migrate

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/fin-agg .
COPY --from=builder /app/migrate .
COPY app-config.yml .
COPY cmd/migrate/migrations ./migrations

CMD ["./fin-agg"]
