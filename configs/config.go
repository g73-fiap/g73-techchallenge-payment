package configs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

type AppConfig struct {
	Port        string
	Environment string

	PaymentBrokerURL string
	NotificationURL  string
	SponsorId        string

	PaymentTable         string
	PaymentTableEndpoint string

	OrderApiUrl      string
	ProductionApiUrl string

	DefaultTimeout int
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GetViperConfig() *viper.Viper {
	return c.viper
}

func (c *Config) ReadConfig() (AppConfig, error) {
	log.Info("Reading Environment Variables")
	c.setupEnvironment()

	log.Info("Setting Config Path")
	c.setupConfigPath()

	log.Info("Reading Config File")
	err := c.viper.ReadInConfig()
	if err != nil {
		return AppConfig{}, fmt.Errorf("error reading config file or env variable, error: %v", err)
	}

	appConfig, err := c.extractConfigVars()
	if err != nil {
		return AppConfig{}, err
	}

	return appConfig, nil
}

func (c *Config) setupConfigPath() {
	// Get proper config directory
	configDirPath := c.viper.GetString("CONFIG_DIR_PATH")
	if configDirPath == "" {
		configDirPath = "./configs"
	}
	log.Infof("ConfigPath %v", configDirPath)
	c.viper.AddConfigPath(configDirPath)
}

func (c *Config) setupEnvironment() {
	c.viper = viper.New()
	c.viper.AutomaticEnv()

	environment := c.viper.GetString("ENVIRONMENT")
	log.Infof("ENVIRONMENT %s", environment)
	c.viper.SetConfigType("yaml")
	c.viper.SetConfigName(environment)
}

func (c *Config) extractConfigVars() (AppConfig, error) {
	appConfig := AppConfig{}

	appConfig.Port = c.viper.GetString("PORT")
	appConfig.Environment = c.viper.GetString("ENVIRONMENT")

	appConfig.PaymentBrokerURL = c.viper.GetString("paymentBroker.url")
	appConfig.NotificationURL = c.viper.GetString("paymentBroker.notificationUrl")
	appConfig.SponsorId = c.viper.GetString("paymentBroker.sponsorId")

	appConfig.PaymentTable = c.viper.GetString("paymentRepository.table")
	appConfig.PaymentTableEndpoint = c.viper.GetString("paymentRepository.endpoint")

	appConfig.OrderApiUrl = c.viper.GetString("ORDER_API_URL")
	appConfig.ProductionApiUrl = c.viper.GetString("PRODUCTION_API_URL")

	appConfig.DefaultTimeout = c.viper.GetInt("DEFAULT_TIMEOUT")

	return appConfig, nil
}
