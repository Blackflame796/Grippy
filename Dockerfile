FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /grippy-server ./cmd/server/main.go

FROM debian:12-slim

WORKDIR /

COPY --from=builder /grippy-server /grippy-server

EXPOSE 8080

CMD ["/grippy-server"]
