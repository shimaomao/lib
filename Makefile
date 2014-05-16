GO=go
.PHONY: all clean test cleanproc

all: test

deps:
	$(GO) get github.com/speedland/wcg
	$(GO) get github.com/influxdb/influxdb-go

test: deps
	@# TARGET_DIR=./
	./scripts/run_tests.sh
