FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Disable CGO for static compilation since we switched to modernc.org/sqlite
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o windmist ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates git

WORKDIR /workspace
COPY --from=builder /app/windmist /usr/local/bin/

ENTRYPOINT ["windmist"]
