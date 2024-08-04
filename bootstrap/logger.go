package bootstrap

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

func (c *Container) initLogger() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	c.logrus = logger.WithContext(context.Background())
}

func (c *Container) Logger() *logrus.Entry {
	return c.logrus
}

func (c *Container) UpdateLogger(updatedLogger *logrus.Entry) {
	c.logrus = updatedLogger
}
