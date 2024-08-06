package main

type Config struct {
	App struct {
		URL string `toml:"url"`
	} `toml:"app"`
	Logger struct {
		Level string `toml:"level"`
	} `toml:"logger"`
	Datasource struct {
		Type string `toml:"type"`
		SQL  struct {
			Host          string `toml:"host"`
			Port          string `toml:"port"`
			Username      string `toml:"user"`
			Password      string `toml:"password"`
			Name          string `toml:"db-name"`
			Ssl           string `toml:"ssl"`
			MigrationsDir string `toml:"migrations-dir"`
		} `toml:"sql"`
	} `toml:"datasource"`
}
