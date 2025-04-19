package types

type RedisConnectionConfig struct {
	Addr       string
	Password   string
	DB         int
	PoolSize   int
	MaxRetries int
}
