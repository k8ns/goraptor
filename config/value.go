package config

import (
	"encoding/base64"
	"sync"
)

type Value struct {
	From      string   `yaml:"from"`
	Value     string   `yaml:"value"`
	Modifiers []string `yaml:"modifiers"`
	o         sync.Once
	v         []byte
}

type Loader interface {
	Load(string) ([]byte, error)
}

type Loaders map[string]Loader

type LoaderFunc func(string) ([]byte, error)

func (f LoaderFunc) Load(s string) ([]byte, error) {
	return f(s)
}

func (cfg *Value) Init(loaders Loaders) (err error) {
	cfg.o.Do(func() {
		cfg.v = []byte(cfg.Value)

		if l, ok := loaders[cfg.From]; ok {
			cfg.v, err = l.Load(cfg.Value)
		}

		for _, m := range cfg.Modifiers {
			if m == "base64" {
				cfg.v, err = base64.StdEncoding.DecodeString(string(cfg.v))
			}
		}
	})

	return err
}

func (cfg *Value) Bytes() []byte {
	return cfg.v
}

func (cfg *Value) String() string {
	return string(cfg.Bytes())
}
