package result_gen

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GenerateResult(filePath string) error {
	err := os.WriteFile(filePath+"/result_err.go", []byte(code1), 0644)
	if err != nil {
		logrus.Errorf("Error writing error with codes: %v", err)
		return err
	}

	cmd := exec.Command("goimports", "-w", filePath+"/result_err.go")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/result_err.go")

	err = os.WriteFile(filePath+"/result_ok.go", []byte(code2), 0644)
	if err != nil {
		logrus.Errorf("Error writing error with codes: %v", err)
		return err
	}

	cmd = exec.Command("goimports", "-w", filePath+"/result_ok.go")
	cmd.Run()

	cmd = exec.Command("go get github.com/goccy/go-json")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/result_err.go")

	return nil
}

const code1 = `package result

import (
	"github.com/goccy/go-json"
)

type ResultErr struct {
	Error string ` + "`json:\"error\"`" + `
	Code  int    ` + "`json:\"code\"`" + `
}

func NewResultErr(err error) *ResultErr {
	var code int

	if errCode, errErr := common.ToErrorWithCode(err); errErr != nil {
		code = int(errCode.GetCode())
	}

	return &ResultErr{
		Error: err.Error(),
		Code:  code,
	}
}

func (r *ResultErr) GetJson() ([]byte, error) {
	return json.Marshal(r)
}`
const code2 = `package result

import (
	"time"

	"github.com/goccy/go-json"
)

type ResultOk struct {
	Result interface{} ` + "`json:\"result\"`" + `
	ExecutionTime time.Duration ` + "`json:\"execution_time\"`" + `
}

func NewResultOk(result interface{}, executionTime time.Duration) *ResultOk {
	return &ResultOk{
		Result: result,
		ExecutionTime: executionTime,
	}
}

func (r *ResultOk) GetJson() ([]byte, error) {
	return json.Marshal(r)
}`
