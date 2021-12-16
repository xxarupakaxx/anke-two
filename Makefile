.PHONY: dev
dev:
	docker-compose -f docker/a/docker-compose.yaml up --build

.PHONY: protobuf-go
protobuf-go:
	@rm -rf src/router/protobuf
	@mkdir -p src/router/protobuf
	@protoc -I ./doc/proto --go_out=src/router/protobuf --go_opt=paths=source_relative doc/proto/*.proto

.PHONY: protobuf-doc
protobuf-doc:
	@protoc --doc_out=html,api.html:doc doc/proto/*.proto