PKG_NAME := env-aws-params
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64
VERSION ?= 0.0.0

TARGETS = $(addprefix target/$(PKG_NAME)_,$(PLATFORMS))

GO := go
RM ?= rm

# some macros to parse that platforms
os = $(word 1,$(subst -, , $@))
arch = $(word 2,$(subst -, , $@))
platform = $(word 2,$(subst _, , $@))

all: build

clean:
	@ $(GO) clean
	@ $(RM) -fr target/

deps:
	$(GO) mod download
	$(GO) mod verify


$(PLATFORMS): deps
	GOOS=$(os) GOARCH=$(arch) $(GO) build -ldflags "-w -s -X main.VersionString=$(VERSION)" -o target/$(PKG_NAME)_$@

$(TARGETS):
	make $(platform)

test: deps
	$(GO) test

fmt:
	$(GO) fmt

build: fmt deps test $(TARGETS)

.PHONY: all clean deps test fmt build
