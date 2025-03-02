package config

type EnvConfig struct {
	Type EnvTypeCfg `yaml:"type" env:"ENV_TYPE" env-default:"dev"`
}

type EnvTypeCfg string

const (
	envProd  EnvTypeCfg = "prod"
	envDev   EnvTypeCfg = "dev"
	envLocal EnvTypeCfg = "local"
)

func (e *EnvConfig) IsProd() bool {
	return e.Type == envProd
}

func (e *EnvConfig) IsDev() bool {
	return e.Type == envDev
}

func (e *EnvConfig) IsLocal() bool {
	return e.Type == envLocal
}

// GetType() возращает тип среды приложения
func (e *EnvConfig) GetType() EnvTypeCfg {
	return e.Type
}
