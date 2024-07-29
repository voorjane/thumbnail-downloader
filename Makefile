build_server:
	go build -o ./.bin/server cmd/server/main.go

build_client:
	go build -o ./.bin/client cmd/client/main.go

run_server: build_server
	./.bin/server

run_client: build_client
	./.bin/client $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

test:
	go test

#ci_build: build_server build_client

.DEFAULT_GOAL := run_server
.PHONY: build_server run_server build_client run_client test #ci_build