package result_gen

import (
	"os"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

func GenerateResult(wg *sync.WaitGroup, errCh chan error, filePath string) {
	for shortPath, code := range codeMap {
		wg.Add(1)
		go func(shortPath string, code string) {
			defer wg.Done()
			err := os.WriteFile(filePath+shortPath, []byte(code), 0644)
			if err != nil {
				logrus.Errorf("Error writing result_err: %v", err)
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
	"/result_err.go": code1,
	"/result_ok.go":  code2,
}

var code1 = `package result

import (
	"github.com/goccy/go-json"
)

type ResultErr struct {
	Error string ` + "`json:\"error\"`" + `
	Code  int    ` + "`json:\"code\"`" + `
}

func NewResultErr(err error) *ResultErr {
	var code int

	if errCode, errErr := error_with_codes.ToErrorWithCode(err); errErr != nil {
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

var code2 = `package result

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
