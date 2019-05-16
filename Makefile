export GOPATH := $(shell pwd)

.PHONY: all clean test

PKG= org/coding/generator \
	 org/coding/ec

all:
	@go install -v $(PKG)

clean:
	@rm -rfv ./bin ./pkg

test:
	@go test -v $(PKG)
