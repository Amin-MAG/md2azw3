FROM --platform=$BUILDPLATFORM golang:1.24.9 AS builder

WORKDIR /go/src/md2awz3

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
    -ldflags "-s -w -X github.com/Amin-MAG/md2awz3/config.AppVersion=${VERSION}" \
    -o /out/md2awz3 ./cmd/md2awz3

########### Main Image ##########
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        wget \
        xz-utils \
        python3 \
        libgl1 \
        libxcb1 \
        libxcb-cursor0 \
        libegl1 \
        libfontconfig1 \
        libfreetype6 \
        libglib2.0-0 \
        libdbus-1-3 \
        libopengl0 && \
    rm -rf /var/lib/apt/lists/*

# Install Calibre and aggressively clean up in a single layer
RUN wget -nv -O- https://download.calibre-ebook.com/linux-installer.sh | sh /dev/stdin install_dir=/opt && \
    # Remove GUI apps
    rm -rf /opt/calibre/bin/calibre \
           /opt/calibre/bin/calibre-server \
           /opt/calibre/bin/calibre-debug \
           /opt/calibre/bin/calibre-smtp \
           /opt/calibre/bin/calibre-complete \
           /opt/calibre/bin/lrf* \
           /opt/calibre/bin/web2disk \
           /opt/calibre/bin/fetch-ebook-metadata \
    # Remove GUI/server Python modules
           /opt/calibre/lib/calibre/gui2 \
           /opt/calibre/lib/calibre/db \
           /opt/calibre/lib/calibre/edit \
           /opt/calibre/lib/calibre/srv \
           /opt/calibre/lib/calibre/live \
    # Remove Qt, Wayland, XCB, ICU libs
           /opt/calibre/lib/libQt*.so* \
           /opt/calibre/lib/libwayland*.so* \
           /opt/calibre/lib/libxcb*.so* \
           /opt/calibre/lib/libxkb*.so* \
           /opt/calibre/lib/libicu*.so* \
           /opt/calibre/lib/libssl*.so* \
           /opt/calibre/lib/libcrypto*.so* \
    # Remove resources
           /opt/calibre/resources/images \
           /opt/calibre/resources/icons \
           /opt/calibre/resources/content-server \
    # Remove build deps
    && apt-get purge -y wget xz-utils && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

ENV PATH="/opt/calibre/bin:${PATH}"

COPY --from=builder /out/md2awz3 /app/md2awz3

ENTRYPOINT ["./md2awz3"]