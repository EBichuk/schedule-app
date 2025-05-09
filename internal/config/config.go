package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db               DbConfig        `yaml:"db"`
	MedicationPeriod MedPeriodConfig `yaml:"medication_period"`
	Logs             LogConfig       `yaml:"logs"`
}

type DbConfig struct {
	Username string `json:"username"`
	Host     string `json:"host:"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	Sslmodel string `json:"sslmode"`
	Password string `json:"password"`
}

type MedPeriodConfig struct {
	Period string `json:"period:"`
	Start  string `json:"start:"`
	End    string `json:"end:"`
}

type LogConfig struct {
	Logfile string `json:"logfile:"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		err := cleanenv.ReadConfig("config.yaml", instance)
		if err != nil {
			disc, _ := cleanenv.GetDescription(instance, nil)
			log.Fatal(disc)
		}
	})
	return instance
}
