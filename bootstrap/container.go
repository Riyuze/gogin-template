package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Container struct {
	ctx    context.Context
	dbw    *sqlx.DB
	dbr    *sqlx.DB
	pgw    *pgxpool.Pool
	pgr    *pgxpool.Pool
	trace  *sdktrace.TracerProvider
	logrus *logrus.Entry
	vip    *viper.Viper
}

func Init() *Container {
	c := &Container{
		ctx: context.Background(),
	}

	c.initLogger()

	c.logrus.Debug("initialized config")
	c.initConfig()

	c.logrus.Debug("initalized telemetry")
	c.initTracer()

	return c
}

func (c *Container) initConfig() {
	vip := viper.New()

	vip.AddConfigPath(".")

	err := vip.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			c.logrus.Info("config file is not found in directory")
		} else {
			c.logrus.Fatal(err)
		}
	}

	c.vip = vip
}

func (c *Container) Close() error {
	c.logrus.Info("closing telemetry connection")
	if c.trace != nil {
		c.trace.Shutdown(c.ctx)
	}

	return nil
}

func (c *Container) GetConfig() *viper.Viper {
	if c.vip == nil {
		c.initConfig()
	}

	return c.vip
}

func (c *Container) CopyStruct(from any, to any) {
	copier.Copy(to, from)
}
