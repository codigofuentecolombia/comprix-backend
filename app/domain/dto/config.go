package dto

import (
	"fmt"

	"gorm.io/gorm"
)

type Config struct {
	GormDB   *gorm.DB
	Settings ConfigSettings
}
type ConfigSettings struct {
	Auth      ConfigSettingsAuth      `yaml:"auth"`
	Email     ConfigSettingsEmail     `yaml:"email"`
	Paths     ConfigSettingsPaths     `yaml:"paths"`
	Server    ConfigSettingsServer    `yaml:"server"`
	Storage   ConfigSettingsStorage   `yaml:"storage"`
	Database  ConfigSettingsDatabase  `yaml:"database"`
	Scrapping ConfigSettingsScrapping `yaml:"scrapping"`
}

type ConfigSettingsEmail struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ConfigSettingsAuth struct {
	Google   ConfigSettingsAuthCredentials `yaml:"google"`
	Facebook ConfigSettingsAuthCredentials `yaml:"facebook"`
}

type ConfigSettingsStorage struct {
	Path string `yaml:"path"`
}

type ConfigSettingsAuthCredentials struct {
	Client   string `yaml:"client"`
	Secret   string `yaml:"secret"`
	Callback string `yaml:"callback"`
}

type ConfigSettingsServer struct {
	Sk     string `yaml:"sk"`
	Port   string `yaml:"port"`
	Host   string `yaml:"host"`
	Debug  bool   `yaml:"debug"`
	IsProd bool   `yaml:"is_prod"`
}

type ConfigSettingsPaths struct {
	Logs   string `yaml:"logs"`
	Assets string `yaml:"assets"`
}
type ConfigSettingsDatabase struct {
	Dsn   string `yaml:"dsn"`
	Debug bool   `yaml:"debug"`
}
type ConfigSettingsScrapping struct {
	TotalTries   int `yaml:"total_tries"`
	MaxGoRutines int `yaml:"max_go_rutines"`
}

func (cnf ConfigSettingsServer) GetAddress() string {
	return fmt.Sprintf("%s:%s", cnf.Host, cnf.Port)
}
