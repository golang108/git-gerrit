tag := $(shell git describe --exact-match --tags 2>/dev/null)
commit := $(shell git rev-parse --short=8 HEAD)
version:=1.0.0

INTERNAL_PKG=github.com/golang108/git-gerrit/cmd

LDFLAGS := $(LDFLAGS) -X $(INTERNAL_PKG).Commit=$(commit)

ifneq ($(tag),)
	LDFLAGS += -X $(INTERNAL_PKG).Version=$(version)
else
	LDFLAGS += -X $(INTERNAL_PKG).Version=$(version)-$(commit)
endif

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









