package grpc

type Config struct {
	Debug          bool `env:"DEBUG" env-default:"true"`
	MetricsPort    int  `env:"METRICS_SERVER_PORT" env-default:"8989"`
	GRPCServerPort int  `env:"GRPC_SERVER_PORT" env-default:"50051"`
}
