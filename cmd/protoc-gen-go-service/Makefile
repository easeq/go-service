gen:
	protoc --proto_path=. example/*.proto \
		--go_out=plugins=grpc:example \
		--go-service_out=logtostderr=true:example; \

gen-options:
	protoc --proto_path=options options/*.proto \
		--go_out=plugins=grpc:options \
		--go_opt=paths=source_relative \

.PHONY: gen