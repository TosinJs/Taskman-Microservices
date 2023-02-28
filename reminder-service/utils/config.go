package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	PORT        string `validate:"required"`
	MONGODB     string `validate:"required"`
	MONGO_URI   string `validate:"required"`
	JWTSECRET   string `validate:"required"`
	RABBITMQURI string `validate:"required"`
}

func LoadConfig(path, name, configType string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	var C Config

	err := viper.ReadInConfig()
	if err != nil {
		return C, err
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		return C, err
	}

	validate := validator.New()
	err = validate.Struct(&C)

	return C, err
}
