package utils

import (
	"encoding/json"

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

type FirebaseConfig struct {
	Type                        string `json:"type" validate:"required"`
	Project_id                  string `json:"project_id " validate:"required"`
	Private_key_id              string `json:"private_key_id" validate:"required"`
	Private_key                 string `json:"private_key" validate:"required"`
	Client_email                string `json:"client_email" validate:"required"`
	Client_id                   string `json:"client_id" validate:"required"`
	Auth_uri                    string `json:"auth_uri" validate:"required"`
	Token_uri                   string `json:"token_uri" validate:"required"`
	Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url" validate:"required"`
	Client_x509_cert_url        string `json:"client_x509_cert_url" validate:"required"`
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

func LoadFirebaseCreds(path, name, configType string) ([]byte, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	var F FirebaseConfig

	err := viper.ReadInConfig()
	if err != nil {
		return []byte{}, err
	}

	err = viper.Unmarshal(&F)
	if err != nil {
		return []byte{}, err
	}

	validate := validator.New()
	err = validate.Struct(&F)

	return json.Marshal(F)
}
