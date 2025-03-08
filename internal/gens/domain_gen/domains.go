package domain_gen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const codeTemplate = `
package {{.PackageName}}

type {{.TypeName}} struct {
	{{range .Fields}}
	{{.Name}} {{.Type}} ` + "`" + `{{.Tag}}` + "`" + `
	{{end}}
}

`

type DomainGenerator struct {
	Files []*Domain
}

type Domain struct {
	FilePath string
	FileName string
	FileExt  string
	Fields   []*Field
}

type Field struct {
	Name string
	Type string
	Tag  string
}

type Code struct {
	PackageName string
	TypeName    string
	Fields      []*Field
}

func NewDomainGenerator(files []*Domain) *DomainGenerator {
	return &DomainGenerator{Files: files}
}

func NewDomain(filePath, fileName, fileExt string, fields []*Field) *Domain {
	return &Domain{FilePath: filePath, FileName: fileName, FileExt: fileExt, Fields: fields}
}

func (c *DomainGenerator) Generate(wg *sync.WaitGroup, errCh chan error) {
	for _, file := range c.Files {
		wg.Add(1)
		go func(file *Domain) {
			defer wg.Done()
			err := generateCode(file)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"file":  file.FileName,
					"error": err,
				}).Error("failed to generate code")
				errCh <- err
				return
			}
		}(file)
	}
}

// Creating file with generated code
func generateCode(file *Domain) error {
	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return err
	}

	packageName := getPackageName(file.FilePath)
	typeName := toCamelCase(file.FileName)

	code := Code{
		PackageName: packageName,
		TypeName:    typeName,
		Fields:      file.Fields,
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
