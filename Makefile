GOFMT_FILES?=$$(find . -name '*.go')

default: test

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@gofmt -w $(GOFMT_FILES)

test: fmtcheck
	@go test -timeout=60s -parallel=4

bin: fmtcheck
	@sh -c "'$(CURDIR)/scripts/build.sh' --release"

dev: fmtcheck
	@sh -c "'$(CURDIR)/scripts/build.sh'"

.NOTPARALLEL:

.PHONY: default fmtcheck fmt test bin dev
