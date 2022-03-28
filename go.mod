module github.com/cyverse-de/job-status-listener

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/cyverse-de/configurate v0.0.0-20190318152107-8f767cb828d9
	github.com/cyverse-de/messaging/v9 v9.1.1
	github.com/cyverse-de/model v0.0.0-20211027151045-62de96618208
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.4.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.31.0
	go.opentelemetry.io/otel v1.6.1
	go.opentelemetry.io/otel/exporters/jaeger v1.6.1
	go.opentelemetry.io/otel/sdk v1.6.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
