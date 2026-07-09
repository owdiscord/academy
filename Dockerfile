ARG GLEAM_VERSION=v1.17.0
ARG GO_VERSION=1.25

# Go build stage
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o academy ./cmd/academy

# Build stage - compile the application
FROM ghcr.io/gleam-lang/gleam:${GLEAM_VERSION}-erlang-alpine AS front_builder

COPY ./frontend /build/frontend

# Build lustre project
RUN cd /build/frontend \
  && rm -rf build/ \
  && gleam run -m lustre/dev build --minify

# Append asset hashes for the frontend so we avoid caching problems
RUN cd /build/frontend && sh append-hash.sh

# Runner stage
FROM scratch

COPY --from=front_builder /build/frontend/dist /frontend
COPY --from=builder /app/academy /academy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 1323
ENTRYPOINT ["/academy"]
