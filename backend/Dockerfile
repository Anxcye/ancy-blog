# File: Dockerfile
# Purpose: Build and package the Go API and migration binaries for production deployment.
# Module: backend deployment, container image layer.
# Related: cmd/server, cmd/migrate, deploy/docker-compose.yml.

FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates gcc musl-dev

ARG TARGETOS=linux
ARG TARGETARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/server ./cmd/server && \
    CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/migrate ./cmd/migrate

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && \
    adduser -S -G app app

COPY --from=builder /out/server /app/server
COPY --from=builder /out/migrate /app/migrate
COPY migrations /app/migrations

USER app

EXPOSE 8080

CMD ["/app/server"]
