# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_all:
	${GOCMD} build -ldflags '-extldflags "-lstdc++"' ./cmd/git-connector

test:
	$(GOTEST) ./...
