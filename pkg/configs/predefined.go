package configs

type HttpServerConfig struct {
	Host string
	Port string
}

type GrpcServerConfig struct {
	Host string
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	Host string
	Port string
}

type LoggerConfig struct {
	Level string
}
