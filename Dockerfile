FROM --platform=$BUILDPLATFORM golang:1.24.9 AS builder

WORKDIR /go/src/md2azw3

ARG TARGETOS
ARG TARGETARCH
ARG BUILDPLATFORM

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION=local

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath \
    -ldflags "-s -w -X github.com/Amin-MAG/md2azw3/config.AppVersion=${VERSION}" \
    -o /out/md2azw3 ./cmd/md2azw3

########### Main Image ##########
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget lib32z1 && \
    wget -q https://archive.org/download/kindlegen_linux_2_6/kindlegen_linux_2.6_i386_v2_9.tar.gz && \
    tar xzf kindlegen_linux_2.6_i386_v2_9.tar.gz -C /usr/local/bin kindlegen && \
    rm kindlegen_linux_2.6_i386_v2_9.tar.gz && \
    apt-get purge -y wget && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/md2azw3 /app/md2azw3

ENTRYPOINT ["./md2azw3"]
