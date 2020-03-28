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
	Muni       IsMuniConfig `json:"muni" yaml:"muni"`
	Parser     string       `json:"parser" yaml:"parser"`
	ResultsDir string       `json:"cache" yaml:"results_dir"`
}

//IsMuniConfig - Is muni config
type IsMuniConfig struct {
	URL       string `json:"url" yaml:"url"`
	Token     string `json:"token" yaml:"token"`
	Course    string `json:"course" yaml:"course"`
	FacultyID int    `json:"faculty_id" yaml:"faculty_id"`
}

type CacheConfig struct {
	Directory string `json:"dir" yaml:"dir"`
	Disabled  bool   `json:"disabled" yaml:"disabled"`
}

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

	return path.Join(appConfigDir, "config.yml"), nil
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
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find configDir directory.
		appConfigDir, err := GetAppConfigDir()
		if err != nil {
			return err
		}

		if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
			createDefaultConfig(appConfigDir)
		}

		filePath := path.Join(appConfigDir, "config.yml")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			createDefaultConfigFile(filePath)
		}

		viper.AddConfigPath(appConfigDir)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Info("Using config file")
	} else {
		log.WithField("file", viper.ConfigFileUsed()).WithError(err).Error("Unable to use the config file")
		return err
	}
	return nil
}

func createDefaultConfig(dir string) error {
	log.WithField("dir", dir).Info("Creating default config since it does not exists!")
	if err := os.Mkdir(dir, 0755); err != nil {
		log.WithError(err).WithField("dir", dir).Error("Unable to create directory")
		return err
	}

	err := createDefaultConfigFile(dir)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultConfigFile(filePath string) error {
	config := Config{
		Muni: IsMuniConfig{
			URL:       "https://is.muni.cz",
			Token:     "",
			Course:    "PB071",
			FacultyID: 1433,
		},
	}
	log.WithField("cfgFile", filePath).Info("Creating config file")

	if err := config.Save(filePath); err != nil {
		return err
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

	if config.ResultsDir == "" {
		var err error
		config.ResultsDir, err = os.Getwd()
		if err != nil {
			log.WithError(err).Warning("Unable to get current working directory")
			return config, err
		}
	}

	return config, nil
}
