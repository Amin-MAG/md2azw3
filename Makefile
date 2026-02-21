IMAGE_NAME = aminmag/md2awz3
PLATFORMS ?= linux/amd64,linux/arm64
VERSION := $(shell cat VERSION)

.PHONY: all
all: build push clean

.PHONY: build
build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest .

.PHONY: push
push:
	@echo ">> Building and pushing $(IMAGE_NAME):$(VERSION) and :latest for $(PLATFORMS)"
	docker buildx build \
		--platform $(PLATFORMS) \
		--build-arg VERSION=$(VERSION) \
		-t $(IMAGE_NAME):$(VERSION) \
		-t $(IMAGE_NAME):latest \
		--push .

.PHONY: clean
clean:
	docker rmi $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest || true

fmt:
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Formatting..."
	@find . -name '*.go' ! -path './vendor/*' -type f -exec go fmt {} \;
	@find . -name '*.go' ! -path './vendor/*' -type f | xargs goimports -w

check:
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Executing go sec checks..."
	@gosec -exclude-dir=bin ./...
	@echo "Executing golangci-lint checks..."
	$(MAKE) lint

lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run --timeout 5m

lint-fix:
	docker pull golangci/golangci-lint:latest && docker run -it --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run --fix --timeout 5m

.PHONY: dev
dev:
	@which air > /dev/null || go install github.com/air-verse/air@latest
	air

.PHONY: swag
swag:
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/md2awz/main.go -o api --parseDependency --parseInternal

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...
