package apf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
)

// Configuration Right now just support TOML
type Configuration struct {
	Logger  []map[string]interface{} `toml:"Logger"`
	Service map[string]interface{}   `toml:"Service"`
	Other   map[string]interface{}   `toml:"Other"`
}

type Configurator func(app *Application)

func (p *Application) WithTOMLConfiguration(f string) *Application {
	return p.WithAppendConf(TOML(f))
}

func TOML(filename string) *Configuration {
	c := DefaultConfiguration()

	tomlAbsPath, err := filepath.Abs(filename)
	if err != nil {
		panic(fmt.Errorf("toml: %w", err))
	}

	data, err := ioutil.ReadFile(tomlAbsPath)
	if err != nil {
		panic(fmt.Errorf("toml: %w", err))
	}

	if _, err := toml.Decode(string(data), &c); err != nil {
		panic(fmt.Errorf("toml: %w", err))
	}

	return c
}

func DefaultConfiguration() *Configuration {
	return &Configuration{
		Logger:  nil,
		Service: nil,
		Other:   nil,
	}
}

func (p *Application) WithAppendConf(c *Configuration) *Application {
	if c == nil {
		return p
	}

	if p.conf == nil {
		p.conf = DefaultConfiguration()
	}

	if c.Service != nil && p.conf.Service == nil {
		p.conf.Service = c.Service
	}

	if c.Logger != nil && p.conf.Logger == nil {
		p.conf.Logger = c.Logger
	}

	if v := c.Other; len(v) > 0 {
		if p.conf.Other == nil {
			p.conf.Other = make(map[string]interface{}, len(v))
		}
		for key, value := range v {
			p.conf.Other[key] = value
		}
	}

	return p
}
