include $(shell git rev-parse --show-toplevel)/include.mk
run:
	export MQX_REDIS_SERVER=127.0.0.1:6379 && \
	export MQX_REDIS_DB=3 && \
	go run .
build:
	sh $(shell git rev-parse --show-toplevel)/scripts/gobuild.sh mqxd $(shell go env GOPATH)/bin/mqxd
build-linux:
	GOOS=linux GOARCH=amd64 sh $(shell git rev-parse --show-toplevel)/scripts/gobuild.sh mqxd $(shell git rev-parse --show-toplevel)/bin/mqxd_linux_amd64