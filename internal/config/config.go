package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

const (
	appName         = "spin"
	configFileName  = "config"
	configFileType  = "yaml"
	configDirectory = ".config/" + appName
)

var ConfigDir string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home, _ = os.Getwd()
	}
	ConfigDir = filepath.Join(home, configDirectory)

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(ConfigDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := ensureConfigDir(); err == nil {
				viper.WriteConfig()
			}
		}
	}
}

func ensureConfigDir() error {
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

func AppDataDir() string {
	if ConfigDir != "" {
		return ConfigDir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home, _ = os.Getwd()
	}
	return filepath.Join(home, configDirectory)
}

func ProfileDataFile() string {
	return filepath.Join(AppDataDir(), "profiles.json")
}

func ActiveProfileFile() string {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return filepath.Join(AppDataDir(), "active_profile")
	}
	return filepath.Join(AppDataDir(), "active_profile.txt")
}

func GetLastFMAPIKey() string {
	return viper.GetString("lastfm_api_key")
}

func GetLastFMAPISecret() string {
	return viper.GetString("lastfm_api_secret")
}

func SetLastFMAPIKey(key string) error {
	viper.Set("lastfm_api_key", key)
	return viper.WriteConfig()
}

func SetLastFMAPISecret(secret string) error {
	viper.Set("lastfm_api_secret", secret)
	return viper.WriteConfig()
}
