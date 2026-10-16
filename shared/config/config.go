package config

type Log struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
	Format string `env:"LOG_FORMAT" envDefault:"console"`
}

type Tracing struct {
	Enabled  bool   `env:"TRACE_ENABLED" envDefault:"false"`
	Endpoint string `env:"TRACE_ENDPOINT" envDefault:"127.0.0.1:6831"`
}

type Metrics struct {
	Enabled bool   `env:"METRICS_ENABLED" envDefault:"false"`
	Addr    string `env:"METRICS_ADDR" envDefault:":9090"`
}

type JWTSecret struct {
	Secret string `env:"JWT_SECRET" envDefault:"dev-secret"`
}
