package bootstrap

import (
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

func (c *Container) Dbw() *sqlx.DB {
	if c.dbw == nil {
		username := c.vip.GetString("database.postgres.write.username")
		password := c.vip.GetString("database.postgres.write.password")

		var err error
		dsn := c.vip.GetString("database.postgres.write.connection")

		// Parse secret if env implemented
		dsn = strings.Replace(dsn, "{username}", username, 1)
		dsn = strings.Replace(dsn, "{password}", password, 1)

		driverName := c.vip.GetString("database.driver")

		if c.trace != nil {
			driverName, err = otelsql.Register(driverName,
				otelsql.AllowRoot(),
				otelsql.TraceQueryWithoutArgs(),
				otelsql.TraceRowsClose(),
				otelsql.TraceRowsAffected(),
				otelsql.WithSystem(semconv.DBSystemPostgreSQL),
			)
			if err != nil {
				c.logrus.Fatal(err)
			}
		}

		conn, err := sqlx.ConnectContext(c.ctx, driverName, dsn)
		if err != nil {
			c.logrus.Panic(err)
		}

		if err := conn.Ping(); err != nil {
			c.logrus.Panic(err)
		}

		conn.SetConnMaxIdleTime(300 * time.Second)

		c.dbw = conn
	}

	return c.dbw
}

func (c *Container) Dbr() *sqlx.DB {
	if c.dbr == nil {
		username := c.vip.GetString("database.postgres.read.username")
		password := c.vip.GetString("database.postgres.read.password")

		var err error
		dsn := c.vip.GetString("database.postgres.read.connection")

		// Parse secret if env implemented
		dsn = strings.Replace(dsn, "{username}", username, 1)
		dsn = strings.Replace(dsn, "{password}", password, 1)

		driverName := c.vip.GetString("database.driver")

		if c.trace != nil {
			driverName, err = otelsql.Register(driverName,
				otelsql.AllowRoot(),
				otelsql.TraceQueryWithoutArgs(),
				otelsql.TraceRowsClose(),
				otelsql.TraceRowsAffected(),
				otelsql.WithSystem(semconv.DBSystemPostgreSQL),
			)
			if err != nil {
				c.logrus.Fatal(err)
			}
		}

		conn, err := sqlx.ConnectContext(c.ctx, driverName, dsn)
		if err != nil {
			c.logrus.Panic(err)
		}

		if err := conn.Ping(); err != nil {
			c.logrus.Panic(err)
		}

		conn.SetConnMaxIdleTime(300 * time.Second)

		c.dbr = conn
	}

	return c.dbr
}

func (c *Container) Pgw() *pgxpool.Pool {
	if c.pgw == nil {
		username := c.vip.GetString("database.postgres.write.username")
		password := c.vip.GetString("database.postgres.write.password")

		var err error
		dsn := c.vip.GetString("database.postgres.write.connection")

		// Parse secret if env implemented
		dsn = strings.Replace(dsn, "{username}", username, 1)
		dsn = strings.Replace(dsn, "{password}", password, 1)

		pool, err := pgxpool.New(c.ctx, dsn)
		if err != nil {
			c.logrus.Panic(err)
		}

		if err := pool.Ping(c.ctx); err != nil {
			c.logrus.Panic(err)
		}

		c.pgw = pool
	}

	return c.pgw
}

func (c *Container) Pgr() *pgxpool.Pool {
	if c.pgw == nil {
		username := c.vip.GetString("database.postgres.read.username")
		password := c.vip.GetString("database.postgres.read.password")

		var err error
		dsn := c.vip.GetString("database.postgres.read.connection")

		// Parse secret if env implemented
		dsn = strings.Replace(dsn, "{username}", username, 1)
		dsn = strings.Replace(dsn, "{password}", password, 1)

		pool, err := pgxpool.New(c.ctx, dsn)
		if err != nil {
			c.logrus.Panic(err)
		}

		if err := pool.Ping(c.ctx); err != nil {
			c.logrus.Panic(err)
		}

		c.pgr = pool
	}

	return c.pgr
}
