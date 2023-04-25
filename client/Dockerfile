FROM golang:1.19-alpine AS builder
ENV CGO_ENABLED=0

WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.* ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=secret,id=application.yml,dst=./application.yml \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

FROM alpine

# Provide the target platform as an environment variable to the application
ARG TARGETPLATFORM
ENV TARGETPLATFORM=${TARGETPLATFORM}

COPY --from=builder /app/bin/service /

CMD /service

