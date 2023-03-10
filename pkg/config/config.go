package config

import (
	"github.com/spf13/viper"
)

// Config represents the application's configuration.
type Config struct {
	MongoURI     string `mapstructure:"mongo_uri"`
	MongoDBName  string `mapstructure:"mongo_db_name"`
	DiscordToken string `mapstructure:"discord_token"`
	GuildID      string `mapstructure:"guild_id"`
	AttendanceID string `mapstructure:"attendance_id"`
}

// LoadConfig loads the application's configuration from the config file.
func LoadConfig(stage string) (*Config, error) {
	switch stage {
	case "dev":
		viper.SetConfigName("config.dev")
	default:
		viper.SetConfigName("config")
	}
	// viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./pkg/config/")
	viper.SetDefault("mongo_uri", "mongodb://localhost:27017")
	viper.SetDefault("discord_token", "my-discord-token")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
