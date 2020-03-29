package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// Config - Application config
type Config struct {
	Muni    MuniConfig `json:"muni" yaml:"muni" mapstructure:"muni"`
	Parser  string     `json:"parser" yaml:"parser" mapstructure:"parser"`
	Results string     `json:"cache" yaml:"results" mapstructure:"results"`
	DryRun  bool       `json:"dryrun" yaml:"dryrun" mapstructure:"dryrun"`
}

//MuniConfig - Is muni config
type MuniConfig struct {
	URL     string `json:"url" yaml:"url" mapstructure:"url"`
	Token   string `json:"token" yaml:"token" mapstructure:"token"`
	Course  string `json:"course" yaml:"course" mapstructure:"course"`
	Faculty int    `json:"faculty_id" yaml:"faculty" mapstructure:"faculty"`
}

const IsStatConfigName = "isstat-config"

// Gets the application configuration directory
func GetAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := path.Join(configDir, "isstat")
	return appConfigDir, nil
}

// GetConfigFilePath - gets a default config file path
func GetConfigFilePath() (string, error) {
	appConfigDir, err := GetAppConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(appConfigDir, IsStatConfigName), nil
}

// Save the config to the specified file
func (config *Config) Save(file string) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		return err
	}
	return err
}

func (config *Config) Dump() (string, error) {
	content, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// SaveToDefaultLocation - Saves a config to the default location ~/.config/isstat/config.yml
func (config *Config) SaveToDefaultLocation() error {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	return config.Save(filePath)
}

// LoadConfig - Loads a config from the configuration
func LoadConfig(cfgFile string) error {
	setDefaults()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find configDir directory.
		appConfigDir, err := GetAppConfigDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(appConfigDir)
		workingDirectory, err := os.Getwd()
		if err == nil {
			viper.AddConfigPath(workingDirectory)
		}
		viper.SetConfigName(IsStatConfigName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Info("Using config file")
	} else {
		log.WithField("file", viper.ConfigFileUsed()).WithError(err).Debug("Config file not found")
		return nil
	}
	return nil
}

// GetAppConfig - Unmarshal the app configuration using the viper
func GetAppConfig() (Config, error) {
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.WithError(err).WithField("file", viper.ConfigFileUsed()).Error("Unable to parse config")
		return config, err
	}

	if config.Results == "" {
		var err error
		config.Results, err = os.Getwd()
		if err != nil {
			log.WithError(err).Warning("Unable to get current working directory")
			return config, err
		}
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("muni.url", "https://is.muni.cz")
	viper.SetDefault("muni.course", "PB071")
	viper.SetDefault("muni.faculty", 1433)
	viper.SetDefault("parser", "default")
	viper.SetDefault("dryrun", false)
}
