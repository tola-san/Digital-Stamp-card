# docker file description about configuration deploy

FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o server ./cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 10000

CMD ["./server"]
