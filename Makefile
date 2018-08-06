GOFMT_FILES?=$$(find . -name '*.go')

default: test

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@gofmt -w $(GOFMT_FILES)

test: fmtcheck
	@go test ./... -timeout=60s -parallel=4

bin: fmtcheck
	@sh -c "'$(CURDIR)/scripts/build.sh' --release"

dev: fmtcheck
	@sh -c "'$(CURDIR)/scripts/build.sh'"

dist: bin
ifndef GH_AUTH_TOKEN
	$(error environment variable GH_AUTH_TOKEN is undefined)
endif
	@sh -c "'$(CURDIR)/scripts/dist.sh' $(GH_AUTH_TOKEN)"

.NOTPARALLEL:

.PHONY: default fmtcheck fmt test bin dev dist
