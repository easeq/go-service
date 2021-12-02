gen:
	protoc --proto_path=. example/*.proto \
		--go_out=example \
		--go-grpc_out=example \
		--go-service_out=logtostderr=true:example; \

gen-options:
	protoc --proto_path=options options/*.proto \
		--go_out=options \
		--go-grpc_out=options \
		--go_opt=paths=source_relative \

.PHONY: gen