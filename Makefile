GO=go
.PHONY: all clean test cleanproc

all: test

deps:
	$(GO) get github.com/speedland/wcg
	$(GO) get github.com/influxdb/influxdb-go
	$(GO) get code.google.com/p/go-html-transform/h5
	$(GO) get code.google.com/p/go-html-transform/css/selector
	$(GO) get code.google.com/p/go.net/html

test: deps
	@# TARGET_DIR=./
	./scripts/run_tests.sh
