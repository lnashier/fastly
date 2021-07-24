package boot

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

// Setup provides app config
func Setup() (*viper.Viper, error) {
	env := "local"
	if psEnv, ok := os.LookupEnv("ENV"); ok {
		env = psEnv
	}
	cfg := viper.New()
	cfg.AddConfigPath("configs/envs/" + env)
	cfg.SetConfigName("app")
	if err := cfg.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "Failed to load config")
	}
	lookupEnv(cfg)
	return cfg, nil
}

func lookupEnv(cfg *viper.Viper) {
	_ = cfg.BindEnv("server.port", "HTTP_PORT")
}
