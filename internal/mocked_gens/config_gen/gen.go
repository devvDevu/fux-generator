package config_gen

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GenerateConfig(filePath string) error {
	err := os.WriteFile(filePath+"/config.go", []byte(code1), 0644)
	if err != nil {
		logrus.Errorf("Error writing error with codes: %v", err)
		return err
	}
	cmd := exec.Command("goimports", "-w", filePath+"/config.go")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/config.go")

	err = os.WriteFile(filePath+"/env.go", []byte(code2), 0644)
	if err != nil {
		logrus.Errorf("Error writing error with codes: %v", err)
		return err
	}
	cmd = exec.Command("goimports", "-w", filePath+"/env.go")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/env.go")

	return nil
}

const code1 = `package config

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	path        string
	envReader   envReader
}

type envReader interface {
	EnvReadConfig(addr string, cfg interface{}) error
}

func MustLoad(ctx context.Context, configPath string, envReader envReader) *Config {
	operation := "config.MustLoad()"

	cfg := new(Config)
	cfg.envReader = envReader

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"config_path": configPath,
		}).WithError(err).Fatal(common.ErrorFailedToFindConfig.SetOperation(operation)) // set operation on custom error
	}

	err = envReader.EnvReadConfig(configPath, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"config_path": configPath,
		}).WithError(err).Fatal(common.ErrorFailedToReadConfig.SetOperation(operation)) // set operation on custom error
	}

	return cfg
}`
const code2 = `
package config

type EnvConfig struct {
	Type EnvTypeCfg ` + "`yaml:\"type\" env:\"ENV_TYPE\" env-default:\"dev\"`" + `
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
}`
