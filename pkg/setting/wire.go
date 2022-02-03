package setting

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"time"
)

// ProvideAppConfig ...
func ProvideAppConfig() *AppConfig {
	config := AppConfig{
		AppName: viper.GetString("APP_NAME"),
		Env:     viper.GetString("APP_ENV"),
	}
	return &config
}

// ProvideHTTPServerConfig ...
func ProvideHTTPServerConfig() *HTTPServerConfig {
	viper.SetDefault("SERVER_READ_TIMEOUT", 5)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 10)
	viper.SetDefault("SERVER_IDLE_TIMEOUT", 120)

	config := &HTTPServerConfig{
		ServerAddr:   viper.GetString("SERVER_ADDR"),
		ReadTimeout:  time.Duration(viper.GetInt("SERVER_READ_TIMEOUT")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("SERVER_WRITE_TIMEOUT")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("SERVER_IDLE_TIMEOUT")) * time.Second,
	}
	return config
}


// DefaultSuperSet ...
var DefaultSuperSet = wire.NewSet(
	ProvideAppConfig,
	ProvideHTTPServerConfig,
)
