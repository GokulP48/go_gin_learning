package config

func ServerPort() string {
	return ":" + AppConfig.Server.Port
}

func DBHost() string {
	return AppConfig.DB.Host
}

func LogLevel() string {
	return AppConfig.Logger.Level
}
