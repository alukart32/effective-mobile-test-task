# syntax=docker/dockerfile:1

########################
# Step 1: Build binary #
########################
FROM golang:alpine as builder

# Add tools.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /build
RUN mkdir ./target

# Fetch dependencies.
COPY ../go.mod ../go.sum ./
RUN go mod download

# Copy src
COPY ../cmd/person/. ./cmd/person/
COPY ../internal/person/. ./internal/person/
COPY ../internal/pkg/. ./internal/pkg/

# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o ./target ./...

###############################
# Step 2: build a small image #
###############################
FROM scratch

WORKDIR /app

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy static executable.
COPY --from=builder /build/target/ ./

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 8080 8081

# Run the binary.
CMD ["/app/person"]