package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io"
	"reflect"
	"strings"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatYaml Format = "yaml"
	FormatEnv  Format = "env"
)

type Embedding struct {
	envPrefix   string
	filePath    string
	configValue any
}

func (e *Embedding) Init(configPtr any, opts ...Option) error {
	for _, opt := range opts {
		opt(e)
	}

	if !(reflect.ValueOf(configPtr).Kind() == reflect.Ptr) {
		panic("expected pointer to config obj")
	}

	if e.envPrefix != "" {
		viper.SetEnvPrefix(e.envPrefix)
	}

	viper.AutomaticEnv()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}

	if e.filePath != "" {
		viper.SetConfigFile(e.filePath)

		err = viper.ReadInConfig()
		if err != nil {
			return err
		}
	}

	err = viper.Unmarshal(&configPtr)
	if err != nil {
		return err
	}

	e.configValue = configPtr

	return nil
}

type Option func(*Embedding)

func WithConfigFile(filePath string) Option {
	return func(p *Embedding) {
		p.filePath = filePath
	}
}

func WithEnvPrefix(envPrefix string) Option {
	return func(p *Embedding) {
		p.envPrefix = envPrefix
	}
}

// Echo prints current config to io.Writer in defined format.
func (e *Embedding) Echo(w io.Writer, format Format) error {
	m := make(map[string]any)
	err := mapstructure.Decode(e, &m)
	if err != nil {
		return err
	}

	switch format {
	case FormatJSON:
		return e.echoJson(m, w)
	case FormatYaml:
		return e.echoYaml(m, w)
	case FormatEnv:
		return e.echoEnv(m, w)
	default:
		return errors.New("unknown format")
	}
}

func (e *Embedding) echoJson(m map[string]any, w io.Writer) error {
	encoder := json.NewEncoder(w)

	encoder.SetIndent("", "\t")
	err := encoder.Encode(m)
	if err != nil {
		return err
	}

	return nil
}

func (e *Embedding) echoYaml(m map[string]any, w io.Writer) error {
	encoder := yaml.NewEncoder(w)

	encoder.SetIndent(2)
	err := encoder.Encode(m)
	if err != nil {
		return err
	}

	return nil
}

func (e *Embedding) echoEnv(m map[string]any, w io.Writer) error {
	queue := make([][]string, 0)

	for key := range m {
		queue = append(queue, []string{key})
	}

	for len(queue) != 0 {
		var (
			path  []string
			value any
		)

		path, queue = queue[0], queue[1:]

		value = m
		for _, key := range path {
			value = value.(map[string]any)[key]
		}

		switch value.(type) {
		case map[string]any:
			for k := range value.(map[string]any) {
				newPath := make([]string, len(path))
				copy(newPath, path)

				newPath = append(newPath, k)
				queue = append(queue, newPath)
			}
		default:
			if e.envPrefix != "" {
				_, err := io.WriteString(w, e.envPrefix)
				if err != nil {
					return err
				}
			}

			for i, key := range path {
				_, err := io.WriteString(w, strings.ToUpper(key))
				if err != nil {
					return err
				}

				if i != len(path)-1 {
					_, err = io.WriteString(w, "_")
					if err != nil {
						return err
					}
				}
			}

			_, err := io.WriteString(w, fmt.Sprintf("=%v\n", value))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
