package config

const (
	CLIENT_MODE = iota
	SERVER_MODE
)

type Config struct {
}

func LoadConfigFromFile(filePath string) Config {

}

func LoadDefaultConfig() Config {

}
