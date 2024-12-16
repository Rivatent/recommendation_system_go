# Указываем базовый образ с Go
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp ./cmd/main.go


FROM gcr.io/distroless/base

COPY --from=builder /app/myapp /myapp

ENTRYPOINT ["/myapp"]