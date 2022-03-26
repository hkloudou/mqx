run:
	export MQX_REDIS_SERVER=127.0.0.1:6379 && \
	export MQX_REDIS_DB=3 && \
	 go run .
build:
	sh ${shell go env GOMODCACHE}/github.com/hkloudou/xlib@v1.0.55/scripts/gobuild.sh mqxd $(shell go env GOPATH)/bin/mqxd
