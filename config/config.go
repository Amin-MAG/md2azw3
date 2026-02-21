package config

import "strings"

var AppVersion = "local"

type Config struct {
	MD2AZW3 struct {
		IsProductionMode bool `env:"IS_PRODUCTION_MODE" env-default:"false" env-description:"Is in production mode"`
		Port             int  `env:"HTTP_PORT" env-default:"8081" env-description:"HTTP server port"`
	}
	Logger struct {
		Level              string `env:"LOGGER_LEVEL" env-default:"debug" env-description:"Log Level for application log"`
		SQLTraceLogEnable  bool   `env:"LOGGER_SQL_TRACE_LOG_ENABLE" env-default:"false" env-description:"Does the log print low level SQL logs"`
		IsReportCallerMode bool   `env:"LOGGER_IS_REPORT_CALLER_MODE" env-default:"false" env-description:"Does the log have report caller"`
		IsPrettyPrint      bool   `env:"LOGGER_IS_PRETTY_PRINT" env-default:"false" env-description:"Pretty JSON Print flag"`
	}
}

// maskString Masks sensitive information with asterisks from string
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

// maskConfig masks the data
func maskConfig(field *string) {
	*field = maskString(*field)
}

// SecureClone creates a secure instance of Config with masking sensitive information
func (c Config) SecureClone() Config {
	sc := c

	// Censor critical values
	//maskConfig(&sc.Spotify.ClientSecret)

	return sc
}
