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
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath \
    -ldflags "-s -w -X github.com/Amin-MAG/md2azw3/config.AppVersion=${VERSION}" \
    -o /out/md2azw3 ./cmd/md2azw3

########### Main Image ##########
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /out/md2azw3 /app/md2azw3

ENTRYPOINT ["./md2azw3"]
