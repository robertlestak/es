package env

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func ReadEnvFiles(f []string) error {
	l := log.WithFields(log.Fields{
		"fn": "ReadEnvFiles",
		"f":  f,
	})
	l.Debug("reading environment file")
	err := godotenv.Load(f...)
	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}
