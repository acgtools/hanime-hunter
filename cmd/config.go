package cmd

type Config struct {
	Log *LogConfig
}

type LogConfig struct {
	Level string
}
