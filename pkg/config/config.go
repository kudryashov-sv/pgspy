package config

import (
	"encoding/json"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configurer interface {
	json.Marshaler

	PostgresURL() string
	NotifyName() string
	OutPath() string
}

type config struct {
	v *viper.Viper
}

func (c *config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.v.AllSettings())
}

func Builder(l *zap.Logger) (Configurer, error) {
	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigName("pgspy")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		l.With(zap.String("component", "config")).Debug("read config error, use defaults", zap.Error(err))
	}

	v.Sub("db").SetDefault("url", "postgres://localhost@pgspy")
	n := v.Sub("notify")
	n.SetDefault("name", "pgspy")
	v.SetDefault("out", "./pgspy.log")

	return &config{v: v}, nil
}

func (c *config) PostgresURL() string {
	return c.v.Sub("db").GetString("url")
}

func (c *config) NotifyName() string {
	return c.v.Sub("notify").GetString("name")
}

func (c *config) OutPath() string {
	return c.v.GetString("out")
}
