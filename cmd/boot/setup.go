package boot

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

// Setup provides app config
func Setup() (*viper.Viper, error) {
	env, ok := os.LookupEnv("ENV")
	if !ok || len(env) < 1 {
		env = "local"
	}
	cfg := viper.New()
	cfg.AddConfigPath("configs/envs/" + env)
	cfg.SetConfigName("app")
	if err := cfg.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "Failed to load config")
	}
	return cfg, nil
}
