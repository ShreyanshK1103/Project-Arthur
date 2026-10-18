# syntax=docker/dockerfile:1

# ---------- build the Go binaries (uses the vendor/ directory, no network needed) ----------
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/api ./cmd/api \
 && CGO_ENABLED=0 go build -trimpath -o /out/worker ./cmd/worker

# ---------- goose, for database migrations ----------
FROM golang:1.26-alpine AS goose
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

# ---------- migrate: applies sql/schema then exits ----------
FROM alpine:3.21 AS migrate
COPY --from=goose /go/bin/goose /usr/local/bin/goose
COPY sql/schema /migrations
CMD ["sh", "-c", "goose -dir /migrations postgres \"$DB_URL\" up"]

# ---------- api ----------
FROM alpine:3.21 AS api
RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 app
COPY --from=build /out/api /usr/local/bin/api
USER app
EXPOSE 8080
CMD ["api"]

# ---------- worker: needs git and the docker CLI to run sandboxed builds ----------
FROM alpine:3.21 AS worker
RUN apk add --no-cache ca-certificates git docker-cli
COPY --from=build /out/worker /usr/local/bin/worker
CMD ["worker"]
