package core

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Path             string `toml:"-"`
	Listen           string
	DataDir          string
	Debug            bool // enable debug mode
	ProxyURL         string
	MidJourneyConfig MidJourneyConfig
	WeChatConfig     WeChatConfig
	SdConfig         SdConfig
	RedisConfig      RedisConfig
	CallbackToken    string
}

type MidJourneyConfig struct {
	Enabled     bool
	UserToken   string
	BotToken    string
	GuildId     string // Server ID
	ChanelId    string // Chanel ID
	CallbackUrl string
}

type WeChatConfig struct {
	Enabled     bool
	CallbackUrl string
}

// SdConfig stable diffusion API config
type SdConfig struct {
	Enabled bool
	ApiURL  string
	ApiKey  string
}

type RedisConfig struct {
	Address  string
	Password string
	Db       int
}

func NewDefaultConfig() *Config {
	return &Config{
		Listen:           "0.0.0.0:9001",
		DataDir:          "./data",
		Debug:            true,
		WeChatConfig:     WeChatConfig{Enabled: true},
		MidJourneyConfig: MidJourneyConfig{Enabled: true},
		RedisConfig: RedisConfig{
			Address:  "localhost:6379",
			Password: "",
			Db:       0,
		},
	}
}

func LoadConfig(configFile string) (*Config, error) {
	var config *Config
	_, err := os.Stat(configFile)
	if err != nil {
		logger.Info("creating new config file: ", configFile)
		config = NewDefaultConfig()
		config.Path = configFile
		// save config
		err := SaveConfig(config)
		if err != nil {
			return nil, err
		}

		return config, nil
	}
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}

	return config, err
}

func SaveConfig(config *Config) error {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(&config); err != nil {
		return err
	}

	return os.WriteFile(config.Path, buf.Bytes(), 0644)
}
