# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINST=$(GOCMD) install

#Binary Name
BINARY_NAME=main

SRC_CLI=$(shell find ./cmd/cli -name *.go)

# Build CLI
build-cli: $(SRC_CLI)
	@$(GOBUILD) -o ./build/$(BINARY_NAME) ./cmd/cli
	@echo "📦 Build CLI Done"

aws-deploy: build-cli
	@./build/main aws deploy
	@echo "🚀 Deploying App to AWS Done"

# Test
test:
	@$(GOTEST) -cover -v ./...
	@echo "🧪 Test Completed"

# Run
run:
	@echo "🚀 Running App"
	@./$(BINARY_NAME)

# Generate Mocks
generate-mocks:
	@$(GOINST) github.com/golang/mock/mockgen@v1.6.0
	@./scripts/generate-mocks.sh

