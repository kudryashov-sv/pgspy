package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

func Builder() (l *zap.Logger, err error) {
	if env := strings.ToLower(os.Getenv("DEVENV")); env == "t" || env == "true" || env == "1" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
