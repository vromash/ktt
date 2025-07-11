FROM golang:1.24.4-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -C cmd -o ../fin-agg

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/fin-agg .
COPY app-config.yml .

EXPOSE 8080

CMD ["./fin-agg"]
