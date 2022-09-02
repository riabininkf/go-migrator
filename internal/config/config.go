package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config = viper.Viper

func New(envPrefix string, cmd *cobra.Command) (*Config, error) {
	var err error

	cnf := viper.New()

	var configPath string
	if configPath, err = cmd.Flags().GetString("config"); err != nil {
		return nil, fmt.Errorf("can't get \"config\" flag: %w", err)
	}

	if len(configPath) == 0 {
		configPath = os.Getenv(fmt.Sprintf("%s_CONFIG", envPrefix))
	}

	if len(configPath) > 0 {
		if err = applyConfigFile(configPath, cnf); err != nil {
			return nil, fmt.Errorf("can't apply config file: %w", err)
		}
	}

	cnf.SetEnvPrefix(envPrefix)

	return cnf, nil
}

func applyConfigFile(path string, cnf *viper.Viper) error {
	var err error

	var b []byte
	if b, err = ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("can't read config file: %w", err)
	}

	expanded := os.ExpandEnv(string(b))
	cnf.SetConfigType("json")

	if err = cnf.ReadConfig(bytes.NewBuffer([]byte(expanded))); err != nil {
		return fmt.Errorf("can't read expanded file: %w", err)
	}

	return nil
}
