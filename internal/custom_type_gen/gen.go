package custom_type_gen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const codeTemplate = `
package {{.PackageName}}

type {{.TypeName}} {{.Type}}

func (ct {{.TypeName}}) {{.MethodName}}() {{.Type}} {
	return {{.Type}}(ct)
} `

type CustomTypeGenerator struct {
	Files []*CustomType
}

func NewCustomTypeGenerator(files []*CustomType) *CustomTypeGenerator {
	return &CustomTypeGenerator{
		Files: files,
	}
}

// FilePath: path to the file Example: "internal/custom_type_gen"
// FileName: name of the file Example: "custom_type"
// FileExt: extension of the file Example: ".go"
// FileType: type of the file Example: "string || int || bool || float64"
type CustomType struct {
	FilePath string
	FileName string
	FileExt  string
	FileType string
}

func NewCustomType(filePath, fileName, fileExt, fileType string) *CustomType {
	return &CustomType{
		FilePath: filePath,
		FileName: fileName,
		FileExt:  fileExt,
		FileType: fileType,
	}
}

type Code struct {
	PackageName string
	TypeName    string
	Type        string
	MethodName  string
}

// Generating code for all files
func (c *CustomTypeGenerator) Generate() error {
	for _, file := range c.Files {
		err := generateCode(file)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  file.FileName,
				"error": err,
			}).Error("failed to generate code")
			return err
		}
	}
	return nil
}

// Creating file with generated code
func generateCode(file *CustomType) error {
	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return err
	}

	methodName := getMethodName(file.FileType)
	packageName := getPackageName(file.FilePath)
	typeName := toCamelCase(file.FileName)

	code := Code{
		PackageName: packageName,
		TypeName:    typeName,
		Type:        file.FileType,
		MethodName:  methodName,
	}

	filePath := filepath.Join(file.FilePath, "/", file.FileName+file.FileExt)
	osFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer osFile.Close()

	tmpl.Execute(osFile, code)

	cmd := exec.Command("goimports", "-w", filePath)
	cmd.Run()

	logrus.Infof("generated code for %s", filePath)
	return nil
}

// Getting package name from file path
func getPackageName(filePath string) string {
	words := strings.Split(filePath, "/")
	return words[len(words)-1]
}

func getMethodName(fileType string) string {
	words := strings.Split(fileType, ".")
	return cases.Title(language.Und).String(strings.ToLower(words[len(words)-1]))
}

// Converting file name to camel case
func toCamelCase(s string) string {
	words := strings.Split(s, "_")
	var builder strings.Builder

	for _, word := range words {
		if word == " " {
			continue
		}
		builder.WriteString(cases.Title(language.Und).String(strings.ToLower(word)))
	}
	return builder.String()
}
