# Render builds from the repository root. The API source lives in gin-api/.
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY gin-api/go.mod gin-api/go.sum ./
RUN go mod download

COPY gin-api/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o server ./cmd/api

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /app/server ./server

EXPOSE 10000
CMD ["./server"]
