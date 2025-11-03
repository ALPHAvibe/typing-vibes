package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type config struct {
	FolderPath   string
	MinLines     int
	MaxLines     int
	MaxTimeLimit int // seconds, 0 = no limit
}

func loadConfig() config {
	viper.SetConfigName("typing_vibes")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/typing_vibes")
	viper.AddConfigPath(".")

	// Set defaults
	homeDir, _ := os.UserHomeDir()
	viper.SetDefault("folder_path", filepath.Join(homeDir, "code"))
	viper.SetDefault("min_lines", 5)
	viper.SetDefault("max_lines", 50)
	viper.SetDefault("max_time_limit", 30)

	viper.ReadInConfig() // Ignore error if config doesn't exist

	return config{
		FolderPath:   viper.GetString("folder_path"),
		MinLines:     viper.GetInt("min_lines"),
		MaxLines:     viper.GetInt("max_lines"),
		MaxTimeLimit: viper.GetInt("max_time_limit"),
	}
}

func saveConfig(cfg config) error {
	viper.Set("folder_path", cfg.FolderPath)
	viper.Set("min_lines", cfg.MinLines)
	viper.Set("max_lines", cfg.MaxLines)
	viper.Set("max_time_limit", cfg.MaxTimeLimit)

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "typing_vibes")
	os.MkdirAll(configDir, 0755)

	configPath := filepath.Join(configDir, "typing_vibes.yaml")
	return viper.WriteConfigAs(configPath)
}

