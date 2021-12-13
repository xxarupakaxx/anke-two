package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Config struct {
	ClientID        string `mapstructure:"client_id" json:"client_id,omitempty"`               // 本番環境向けのクライアントID (default: "")
	DevClientID     string `mapstructure:"dev_client_id" json:"dev_client_id,omitempty"`       // ローカル開発環境向けのクライアントID (default: "")
	MariaDBHostname string `mapstructure:"mariadb_hostname" json:"mariadb_hostname,omitempty"` // DB のホスト (default: "mariadb")
	MariaDBDatabase string `mapstructure:"mariadb_database" json:"mariadb_database,omitempty"` // DB の DB 名 (default: "anke-two")
	MariaDBUsername string `mapstructure:"mariadb_username" json:"mariadb_username,omitempty"` // DB のユーザー名 (default: "root")
	MariaDBPassword string `mapstructure:"mariadb_password" json:"mariadb_password,omitempty"` // DB のパスワード (default: "password")
}

func GetConfig() (*Config, error) {
	viper.SetDefault("Client_ID", "")
	viper.SetDefault("Dev_Client_ID", "")
	viper.SetDefault("MariaDB_Hostname", "mariadb")
	viper.SetDefault("MariaDB_Database", "anke-two")
	viper.SetDefault("MariaDB_Username", "root")
	viper.SetDefault("MariaDB_Password", "password")

	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("Unable to find config.json, default settings or environmental variables are to be used.")
		} else {
			return nil, fmt.Errorf("Error: failed to load config.json - %s ", err)
		}
	}

	var c *Config

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("Error: failed to parse configs - %s ", err)
	}

	return c, nil
}
