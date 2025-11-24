# BINARY_NAME=NAME

check-quality:
	# make lint
	make fmt
	make vet

lint:
	golangci-lint run 

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

# build:
	# mkdir -p bin
	# GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux cmd/main.go
	# @echo "Build complete"

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: run
run:
	go run cmd/main.go $(RUN_ARGS)
# run: build
	# ./bin/${BINARY_NAME}-linux $(RUN_ARGS)

.PHONY: debug
debug:
	DEBUG=true go run cmd/main.go

clean:
	go clean
	rm -rf bin

test-local:
	make tidy
	gotest -v ./...

test:
	make tidy
	gotest -v ./... -coverprofile=coverage.out -json > report.json

coverage:
	make test
	go tool cover -html=coverage.out

.PHONY: all test and build
all:
	make check-quality
	make test
	make build
