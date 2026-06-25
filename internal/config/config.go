// Package config loads all our env variables into a struct
package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	DatabaseURI     string   `env:"DATABASE_URI"`
	AthenaDBURI     string   `env:"ATHENA_DB_URI"`
	ModmailDBURI    string   `env:"MODMAIL_DB_URI"`
	BotToken        string   `env:"DISCORD_BOT_TOKEN"`
	ClientID        string   `env:"DISCORD_CLIENT_ID"`
	ClientSecret    string   `env:"DISCORD_CLIENT_SECRET"`
	RedirectURI     string   `env:"DISCORD_REDIRECT_URI"`
	PrivateChannels []string `env:"PRIVATE_CHANNELS"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// dereference pointer to get the struct value
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envKey := field.Tag.Get("env")

		// skip fields without an env tag
		if envKey == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			return nil, fmt.Errorf("missing required environment variable: %s", envKey)
		}

		fieldVal := v.Field(i)

		// skip unexported fields
		if !fieldVal.CanSet() {
			continue
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(envVal)
		case reflect.Slice:
			switch fieldVal.Type().Elem().Kind() {
			case reflect.String:
				parts := strings.Split(envVal, ",")
				for i, p := range parts {
					parts[i] = strings.TrimSpace(p)
				}
				fieldVal.Set(reflect.ValueOf(parts))
			default:
				return nil, fmt.Errorf("unsupported slice element type %s for %s", fieldVal.Type().Elem().Kind(), field.Name)
			}
		default:
			return nil, fmt.Errorf("unsupported field type %s for %s", fieldVal.Kind(), field.Name)
		}
	}

	return cfg, nil
}
