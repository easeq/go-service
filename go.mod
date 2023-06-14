module github.com/easeq/go-service

go 1.16

require (
	github.com/Netflix/go-env v0.0.0-20210215222557-e437a7e7f9fb
	github.com/easeq/go-consul-registry/v2 v2.1.0
	github.com/easeq/go-redis-access-control v0.0.6
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gofiber/fiber/v2 v2.43.0
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0
	github.com/hashicorp/consul/api v1.8.1 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/nats-io/nats-server/v2 v2.9.8 // indirect
	github.com/nats-io/nats.go v1.20.0
	github.com/nsqio/go-nsq v1.1.0
	github.com/stretchr/testify v1.8.0
	go.etcd.io/etcd/client/v3 v3.5.6
	go.opentelemetry.io/otel v1.11.1
	go.opentelemetry.io/otel/exporters/jaeger v1.11.1
	go.opentelemetry.io/otel/sdk v1.11.1
	go.opentelemetry.io/otel/trace v1.11.1
	go.uber.org/zap v1.23.0
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
)
