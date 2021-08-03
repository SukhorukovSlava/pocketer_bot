package config

import "github.com/spf13/viper"

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	AuthServerPort    string
	TelegramBotURL    string `mapstructure:"bot_url"`
	DBPath            string `mapstructure:"db_path"`

	Messages Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default         string `mapstructure:"default"`
	InvalidLink     string `mapstructure:"invalid_link"`
	Unauthorized    string `mapstructure:"unauthorized"`
	UnableToAddLink string `mapstructure:"unable_to_add_link"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorize  string `mapstructure:"already_authorize"`
	AddedSuccessfully string `mapstructure:"added_successfully"`
	UnknownCmd        string `mapstructure:"unknown_cmd"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}
	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}
	if err := viper.BindEnv("auth_server_port"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.PocketConsumerKey = viper.GetString("consumer_key")
	cfg.AuthServerURL = viper.GetString("auth_server_url")
	cfg.AuthServerPort = viper.GetString("auth_server_port")

	return nil
}
