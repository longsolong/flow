package setting

import (
	"time"
)

// AppConfig ...
type AppConfig struct {
	AppName    string
	Env        string
}

// HTTPServerConfig ...
type HTTPServerConfig struct {
	ServerAddr string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
}
