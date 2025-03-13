package db

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/k8ns/goraptor/config"
)

type Config struct {
	Engine        string
	Username      config.Value
	Password      config.Value
	Database      string
	Host          config.Value
	Port          string
	Options       map[string]string
	RunMigrations bool `yaml:"run_migrations"`
}

func (c *Config) Init(l config.Loaders) (err error) {
	if err = c.Username.Init(l); err != nil {
		return
	}

	if err = c.Password.Init(l); err != nil {
		return
	}

	if err = c.Host.Init(l); err != nil {
		return
	}

	return
}

func NewDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Engine, cfg.Dsn())
	if err != nil {
		return nil, fmt.Errorf("couldn't open connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't ping: %v", err)
	}

	return db, nil
}

func (c *Config) dsnArgs() []any {
	a := []any{c.Username.String(), c.Password.String(), c.Host.String(), c.Port, c.Database}

	b := bytes.Buffer{}
	for k, v := range c.Options {
		if b.Len() > 0 {
			b.Write([]byte("&"))
		}
		b.Write([]byte(k))
		b.Write([]byte("="))
		b.Write([]byte(v))
	}

	a = append(a, b.String())
	return a
}

func (c *Config) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.dsnArgs()...)
}

func DriverName(conn *sql.DB) string {
	dr := fmt.Sprintf("%#v", conn.Driver())
	return dr[1:6]
}
