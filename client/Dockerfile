FROM golang:1.19-alpine AS builder
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

FROM alpine

ARG TARGETPLATFORM
ENV TARGETPLATFORM=${TARGETPLATFORM}

COPY --from=builder /app/bin/service /

CMD /service

