GOCMD ?= go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean -cache
GOGET = $(GOCMD) get
BINARY_NAME = isp-checker


default: get build


get:
	@go mod tidy


build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v -ldflags="-s -w" ./cmd/blacklist/main.go
	upx $(BINARY_NAME)


clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

