package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/inconshreveable/log15"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"

	"github.com/fatih/structs"
)

// AsEnvVariables sets struct values from environment variables
func AsEnvVariables(o interface{}, prefix string, skipCommented bool) map[string]string {
	r := map[string]string{}
	prefix = strings.ToUpper(prefix)
	delim := "_"
	if prefix == "" {
		delim = ""
	}
	fields := structs.Fields(o)
	for _, f := range fields {
		if skipCommented {
			tag := f.Tag("commented")
			if tag != "" {
				commented, err := strconv.ParseBool(tag)
				log.Error("Unable to parse tag value", "error", err)
				if commented {
					continue
				}
			}
		}
		if structs.IsStruct(f.Value()) {
			rf := AsEnvVariables(f.Value(), prefix+delim+f.Name(), skipCommented)
			for k, v := range rf {
				r[k] = v
			}
		} else {
			r[prefix+"_"+strings.ToUpper(f.Name())] = fmt.Sprintf("%v", f.Value())
		}
	}
	return r
}

// Load a config
func Load(conf interface{}, envPrefix string, cfgFile string) error {
	// Apply defaults first
	defaults.SetDefaults(conf)

	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// Overrides with environment
	for k := range AsEnvVariables(conf, "", false) {
		envName := fmt.Sprintf("%s_%s", upPrefix, k)
		viper.BindEnv(strings.ToLower(strings.Replace(k, "_", ".", -1)), envName)
	}

	// Apply file settings
	if cfgFile != "" {
		// If the config file doesn't exists, let's exit
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return xerrors.Errorf("Unable to open non-existing file '%s': %w", cfgFile, err)
		}

		log.Info("Load settings from file", "path", cfgFile)

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return xerrors.Errorf("Unable to decode config file '%s': %w", cfgFile, err)
		}
	}

	// Update viper values
	if err := viper.Unmarshal(conf); err != nil {
		return xerrors.Errorf("Unable to apply config '%s': %w", cfgFile, err)
	}

	// No error
	return nil
}
