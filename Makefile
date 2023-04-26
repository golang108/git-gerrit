tag := $(shell git describe --always --all --dirty --broken 2>/dev/null)
commit := $(shell git rev-parse --short=8 HEAD)
version:=$(shell git rev-list --count HEAD)

INTERNAL_PKG=github.com/golang108/git-gerrit/cmd

LDFLAGS := $(LDFLAGS) -X $(INTERNAL_PKG).Commit=$(commit)
LDFLAGS += -X $(INTERNAL_PKG).Version=$(version)-$(tag)


all: tidy build


.PHONY: version
version:
	@echo $(version)-$(commit)


.PHONY: build
build:
	GOOS=linux GOARCH=amd64 \
go build -trimpath -ldflags "$(LDFLAGS)" -o git-gerrit .
	GOOS=windows GOARCH=amd64 \
go build -trimpath -ldflags "$(LDFLAGS)" -o git-gerrit.exe .
	GOOS=darwin GOARCH=amd64 \
go build -trimpath -ldflags "$(LDFLAGS)" -o git-gerrit.mac .













.PHONY: tidy
tidy:
	go mod verify
	go mod tidy

.PHONY: clean
clean:
	rm -rf git-gerrit









