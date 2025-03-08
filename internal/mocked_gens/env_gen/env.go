package env_gen

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GenerateEnv(filePath string) error {
	os.MkdirAll(filePath, 0755)
	err := os.WriteFile(filePath+"/env.go", []byte(code1), 0644)
	if err != nil {
		logrus.Errorf("Error writing env: %v", err)
		return err
	}

	cmd := exec.Command("go get github.com/ilyakaznacheev/cleanenv")
	cmd.Run()

	cmd = exec.Command("goimports", "-w", filePath+"/env.go")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/env.go")

	return nil
}

const code1 = `package env

import "github.com/ilyakaznacheev/cleanenv"

type EnvReader struct{}

func NewEnvReader() *EnvReader {
	return &EnvReader{}
}

func (e *EnvReader) EnvReadConfig(configPath string, cfg interface{}) error {
	return cleanenv.ReadConfig(configPath, cfg)
}`
