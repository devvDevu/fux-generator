package env

import "github.com/ilyakaznacheev/cleanenv"

type EnvReader struct{}

func NewEnvReader() *EnvReader {
	return &EnvReader{}
}

func (e *EnvReader) EnvReadConfig(configPath string, cfg interface{}) error {
	return cleanenv.ReadConfig(configPath, cfg)
}
