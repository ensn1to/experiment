package gopool

const (
	defaultScalaThreshold = 1
)

type Config struct {
	ScaleThreshold int32
}

func NewConfig() *Config {
	return &Config{
		ScaleThreshold: defaultScalaThreshold,
	}
}
