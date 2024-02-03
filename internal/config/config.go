package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

var (
	configWatcher  = viper.New()
	configFileName = "config"
	configType     = "yaml"
	configPath     = "./configs"
)

func Get() (*GlobalRateLimiter, error) {
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := Default()
	err = watchFile[GlobalRateLimiter](config, configFileName, configWatcher)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func watchFile[T any](bindTo *T, file string, watcher *viper.Viper) error {
	watcher.AddConfigPath(configPath)     // local folder
	watcher.AddConfigPath("/app/configs") // required for container
	watcher.SetConfigName(file)
	watcher.SetConfigType(configType)
	watcher.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := watcher.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in configs. err: %s", err)
	}

	if err := watcher.Unmarshal(bindTo); err != nil {
		return fmt.Errorf("failed to unmarshal configs. err: %s", err)
	}

	watcher.WatchConfig()
	watcher.OnConfigChange(func(event fsnotify.Event) {
		if err := watcher.Unmarshal(bindTo); err == nil {
			zap.L().Info(fmt.Sprintf("configs %s updated", file))
		} else {
			zap.L().Error(fmt.Sprintf("could not unmarshal %s configs.", file), zap.Error(err))
		}
	})

	return nil
}
