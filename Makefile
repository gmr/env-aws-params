.DEFAULT_GOAL = build

GO := go
DEP := dep

PKG_NAME=env-aws-params

all: build

clean:
	@ $(GO) clean
	@ rm -rf target/

PLATFORMS := linux-amd64 linux-arm64 darwin-amd64

deps: 
	@ $(DEP) ensure
	@ $(DEP) check

os = $(word 1,$(subst -, ,$@))
arch = $(word 2,$(subst -, ,$@))
platform = $(word 2,$(subst _, ,$@))

$(PLATFORMS): deps
	GOOS=$(os) GOARCH=$(arch) $(GO) build \
		-ldflags "-w -s -X main.VersionString=v${TRAVIS_TAG}" \
		-o target/$(PKG_NAME)_$@ 

TARGETS = $(addprefix target/$(PKG_NAME)_,$(PLATFORMS))

$(TARGETS): 
	make $(platform)

test: deps
	$(GO) test

fmt:
	$(GO) fmt

build: fmt deps test $(TARGETS)

.PHONY: deps build