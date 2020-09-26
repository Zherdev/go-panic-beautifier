package cfg

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`

	MainPageUrl      string `yaml:"main_page_url"`
	MainPageTemplate string `yaml:"main_page_template"`

	ReadTimeout time.Duration `yaml:"read_timeout"`
	WriteTimout time.Duration `yaml:"write_timeout"`

	LogLevel uint32 `yaml:"log_level"`
}

func NewConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open config file %q", filename)
	}

	result := &Config{}

	d := yaml.NewDecoder(file)
	err = d.Decode(result)
	if err != nil {
		return result, errors.Wrapf(err, "can't unmarshal config from %q", filename)
	}

	return result, err
}
