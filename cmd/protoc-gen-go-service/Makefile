gen:
	protoc --proto_path=. example/*.proto \
		--go_out=plugins=grpc:example \
		--go-service_out=logtostderr=true:example; \
	protoc --proto_path=options options/*.proto \
		--go_out=plugins=grpc:options \

.PHONY: gen