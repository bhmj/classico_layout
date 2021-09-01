package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var errConfigTypeNotSupported = errors.New("unsupported config file type")

type Pallet struct {
	Large  int `long:"large" env:"LARGE" description:"Number of large pieces in one layer" default:"32"`
	Medium int `long:"medium" env:"MEDIUM" description:"Number of medium pieces in one layer" default:"32"`
	Small  int `long:"small" env:"SMALL" description:"Number of small pieces in one layer" default:"8"`
	Layers int `long:"layers" env:"LAYERS" description:"Number of layers on a pallet" default:"12"`
}

type Color struct {
	Mode  string `long:"mode" env:"MODE" description:"Color mode" default:"random" choice:"random" choice:"gradient"`
	Color map[string]ColorDefinition
}

type ColorDefinition struct {
	Name    string `long:"name" env:"NAME" description:"Color name (applied to img/<name>_*.png)"`
	Pallets int    `long:"pallets" env:"PALLETS" description:"Number of pallets of this color"`
}

type Pavement struct {
	Width int `long:"width" env:"WIDTH" description:"Width of the pavement in small piece numbers" default:"10"`
}

// Config contains all parameters
type Config struct {
	ConfigFile string   `long:"config-file" env:"CLASSICO_CONFIG_FILE" description:"Classico config file path (json and yaml formats supported)"`
	Pallet     Pallet   `env-namespace:"CLASSICO_PALLET" namespace:"pallet" group:"Pallet configuration"`
	Pavement   Pavement `env-namespace:"CLASSICO_ROAD" namespace:"pavement" group:"Pavement configuration"`
	Color      Color    `env-namespace:"CLASSICO_COLOR" namespace:"color" group:"Color configuration"`
	LogLevel   string   `long:"log-level" env:"CLASSICO_LOG_LEVEL" description:"Log level" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"dpanic" choice:"panic" choice:"fatal"` // nolint:staticcheck
}

// New returns instance of config
func New() *Config {
	return &Config{
		LogLevel: "debug",
	}
}

type configType string

const (
	jsonConfig configType = "json"
	yamlConfig configType = "yaml"
)

// ReadFromFile reads config from file
func (t *Config) ReadFromFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()
	configType, err := t.getConfigType(fname)
	if err != nil {
		return err
	}
	return t.parseRead(f, configType)
}

func (t *Config) getConfigType(fname string) (configType, error) {
	ext := filepath.Ext(fname)
	switch ext {
	case ".yaml", ".yml":
		return yamlConfig, nil
	case ".json":
		return jsonConfig, nil
	default:
		return "", fmt.Errorf("%w: %s", errConfigTypeNotSupported, ext)
	}
}

func (t *Config) parseRead(f io.Reader, fileType configType) error {
	// pass secrets through env
	conf, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("config file ReadAll: %w", err)
	}
	rx := regexp.MustCompile(`{{(\w+)}}`)
	for {
		matches := rx.FindSubmatch(conf)
		if matches == nil {
			break
		}
		v := os.Getenv(strings.ToUpper(string(matches[1])))
		conf = bytes.ReplaceAll(conf, matches[0], []byte(v))
	}

	switch fileType {
	case jsonConfig:
		if err := json.Unmarshal(conf, &t); err != nil {
			return fmt.Errorf("config file json.Unmarshal: %w", err)
		}
	case yamlConfig:
		if err := yaml.Unmarshal(conf, &t); err != nil {
			return fmt.Errorf("config file yaml.Unmarshal: %w", err)
		}
	}
	return nil
}
