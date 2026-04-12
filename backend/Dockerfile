FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:3.22
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
RUN mkdir -p uploads && chown -R app:app /app
USER app

EXPOSE 8080
CMD ["./server"]
