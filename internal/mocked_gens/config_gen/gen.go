package config_gen

import (
	"os"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

func GenerateConfig(wg *sync.WaitGroup, errCh chan error, filePath string) {
	for shortPath, code := range codeMap {
		wg.Add(1)
		go func(shortPath string, code string) {
			defer wg.Done()
			err := os.WriteFile(filePath+shortPath, []byte(code), 0644)
			if err != nil {
				logrus.Errorf("Error writing config: %v", err)
				errCh <- err
				return
			}
			cmd := exec.Command("goimports", "-w", filePath+shortPath)
			cmd.Run()

			logrus.Infof("generated code for %s", filePath+shortPath)
		}(shortPath, code)
	}
}

var codeMap = map[string]string{
	"/config.go": code1,
	"/env.go":    code2,
}

var code1 = `package config

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
		}).WithError(err).Fatal(error_with_codes.ErrorFailedToFindConfig.SetOperation(operation)) // set operation on custom error
	}

	err = envReader.EnvReadConfig(configPath, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"config_path": configPath,
		}).WithError(err).Fatal(error_with_codes.ErrorFailedToReadConfig.SetOperation(operation)) // set operation on custom error
	}

	return cfg
}`

var code2 = `
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
