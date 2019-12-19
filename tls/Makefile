.PHONY: all
all: help

GOLANG_IMG=quay.io/hyper/go:1.13.5

.PHONY: help
help:
	#########################################################################
	#
	# protoc      - compile *.proto to golang code
	# server      - start the grpc server
	# server-down - stop the grpc server
	# client      - run the grpc client
	#
	#########################################################################
	@echo "Enjoy!"

.PHONY: protoc
protoc:
	docker run --rm -t -v "${PWD}:/app" -w "/app" ${GOLANG_IMG} \
		bash -c 'make _protoc; gofmt -w proto'

.PHONY: _protoc
_protoc:
	find ./proto -name "*.proto" | xargs \
	/usr/bin/protoc \
		-I proto \
		-I $$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I $$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
		-I $$GOPATH/src \
		--gogo_out=plugins=grpc:./proto/ \
		--grpc-gateway_out=logtostderr=true:./proto/ \
		--swagger_out=./proto/

.PHONY: logs
logs:
	docker logs -f --tail=30 tls-micro-demo-server

.PHONY: server
server: server-down
	docker run -d --name=tls-micro-demo-server --rm \
	-p 18888:18888 -p 19999:19999 \
	-v "${PWD}:/app" -w "/app" ${GOLANG_IMG} go run server/main.go
	make logs

.PHONY: server-down
server-down:
	docker rm -f tls-micro-demo-server || true

.PHONY: client
client:
	docker exec -w "/app" tls-micro-demo-server go run client/main.go
