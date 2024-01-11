package config

import (
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type externalcfg struct {
	Xendit struct {
		API_KEY string `mapstructure:"XENDIT_API_KEY"`
	} `mapstructure:",squash"`
}

type dbcfg struct {
	DSN         string `mapstructure:"DB_DSN"`
	MaxOpenConn int    `mapstructure:"DB_MAX_OPEN_CONN"`
	MaxIdleConn int    `mapstructure:"DB_MAX_IDLE_CONN"`
	MaxIdleTime string `mapstructure:"DB_MAX_IDLE_TIME"`
}

type jwtcfg struct {
	Secret string `mapstructure:"JWT_SECRET"`
}

type limitercfg struct {
	Enabled bool    `mapstructure:"LIMITER_ENABLED"`
	Rps     float64 `mapstructure:"LIMITER_RPS"`
	Burst   int     `mapstructure:"LIMITER_BURST"`
}

type smtpcfg struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     int    `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
	Sender   string `mapstructure:"SMTP_SENDER"`
}

type Config struct {
	Port        int         `mapstructure:"API_PORT"`
	Environment string      `mapstructure:"ENVIRONMENT"`
	Db          dbcfg       `mapstructure:",squash"`
	Jwt         jwtcfg      `mapstructure:",squash"`
	Limiter     limitercfg  `mapstructure:",squash"`
	Smtp        smtpcfg     `mapstructure:",squash"`
	External    externalcfg `mapstructure:",squash"`

	// cors struct {
	// 	trustedOrigins []string
	// }
}

func New() (*Config, error) {
	cfg := &Config{}

	viper.SetDefault("API_PORT", 8080)
	viper.SetDefault("ENVIRONMENT", "development")

	viper.SetDefault("DB_DSN", os.Getenv("DB_DSN"))
	viper.SetDefault("DB_MAX_OPEN_CONN", 25)
	viper.SetDefault("DB_MAX_IDLE_CONN", 25)
	viper.SetDefault("DB_MAX_IDLE_TIME", "15m")

	viper.SetDefault("LIMITER_ENABLED", true)
	viper.SetDefault("LIMITER_RPS", 2)
	viper.SetDefault("LIMITER_BURST", 4)

	viper.SetDefault("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	viper.SetDefault("SMTP_PORT", 2525)
	viper.SetDefault("SMTP_USERNAME", "0f3dffb73b2846")
	viper.SetDefault("SMTP_PASSWORD", "6c22cf79f49adc")
	viper.SetDefault("SMTP_SENDER", "Pixelrental <no-reply@pixelrental.support.net>")

	viper.SetDefault("JWT_SECRET", os.Getenv("JWT_SECRET"))
	viper.SetDefault("XENDIT_API_KEY", os.Getenv("XENDIT_API_KEY"))

	viper.AutomaticEnv()
	// viper.SetConfigType("env")
	// viper.SetConfigName(".env")
	// viper.AddConfigPath(".")
	// viper.SetConfigFile("env")

	// if err := viper.ReadInConfig(); err != nil {
	// 	return nil, err
	// }

	if err := viper.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.IgnoreUntaggedFields = true
		dc.ErrorUnused = true
	}); err != nil {
		return nil, err
	}

	if structs.HasZero(&cfg) {
		return nil, fmt.Errorf("config type has zero value")
	}
	return cfg, nil
}
