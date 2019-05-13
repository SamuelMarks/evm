export RM?=rm
export GLIDE?=glide
export GO?=go

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	$(GLIDE) install

# install compiles and places the binary in GOPATH/bin
install:
	$(GO) install \
	 	--ldflags '-extldflags "-static"' \
		./cmd/evm

# build compiles and places the binary in /build
build:
	$(GO) build \
		--ldflags '-extldflags "-static"' \
		-o build/evm ./cmd/evm/

# dist builds binaries for all platforms and packages them for distribution
dist:
	@BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/dist.sh'"

test:
	$(GLIDE) novendor | xargs go test

clean:
	$(GLIDE) cc
	$(RM) -rf vendor glide.lock

.PHONY: vendor install build test clean
