FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/server ./cmd/server/

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY config/ config/
EXPOSE 8080
ENTRYPOINT ["/app/server"]
