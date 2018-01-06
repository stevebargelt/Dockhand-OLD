SOURCEDIR=.
BINARY=bin/dockhand
VET_REPORT = vet.report
TEST_REPORT = tests.xml

# These will be provided to the target
VERSION=0.1.0
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

CURRENT_DIR=$(shell pwd)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-s -w -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

.PHONY:bin
bin:
	go build -o $(BINARY) $(LDFLAGS) main.go

.PHONY:pi
pi:
	env GOOS=linux GOARCH=arm go build -o $(BINARY) $(LDFLAGS) main.go

.PHONY:test
test:
	go test -cover -v -race ./...

.PHONY:vet
vet:
	go vet ./...

.PHONY:docker
docker:
	docker build -t dockhand -f Dockerfile .
# docker run --name dockhand -it dockhand

.PHONY:build
build: clean test bin

.PHONY:clean
clean:
	rm -rf bin
