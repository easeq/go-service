gen:
	protoc --proto_path=example example/*.proto \
		--go_out=plugins=grpc:example-out \
		--go-service_out=logtostderr=true:example-out \

.PHONY: gen clean run