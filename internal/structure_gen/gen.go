package structure_gen

import (
	"ca-generator/internal/custom_type_gen"
	"ca-generator/internal/domain_gen"
	"ca-generator/internal/folders_gen"
	"ca-generator/internal/json_file_gen"
	"ca-generator/internal/mocked_gens/config_gen"
	"ca-generator/internal/mocked_gens/env_gen"
	"ca-generator/internal/mocked_gens/error_with_codes_gen"
	"ca-generator/internal/mocked_gens/result_gen"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type Generator struct {
	Settings       map[string]interface{}
	Folders        []string
	Files          []*custom_type_gen.CustomType
	Domain         []*domain_gen.Domain
	ErrorWithCodes string
	Config         string
	Result         string
	Env            string
}

func NewGenerator(settingsFilePath string) (*Generator, error) {
	// reading json data
	jsonData, err := os.ReadFile(settingsFilePath)
	if err != nil {
		logrus.Errorf("Error reading settings.json: %v", err)
		logrus.Info("Creating example settings.json")
		logrus.Warn("Please edit settings.json and run the program again")
		json_file_gen.GenerateJsonFile()
		return nil, err
	}

	// unmarshalling json data
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {
		logrus.Errorf("Error unmarshalling settings.json: %v", err)
		return nil, err
	}

	g := &Generator{
		Settings: config,
	}

	return g, nil
}

func (g *Generator) Generate() {
	// formating settings data
	formatingSettingsData(g.Settings["settings"].(map[string]interface{}), "", &g.Folders, &g.Files, &g.Domain, &g.ErrorWithCodes, &g.Config, &g.Result, &g.Env)

	// generating structure
	err := folders_gen.NewFolderGenerator(g.Folders).Generate()
	if err != nil {
		logrus.Fatalf("Error generating folders: %v", err)
	}
	err = custom_type_gen.NewCustomTypeGenerator(g.Files).Generate()
	if err != nil {
		logrus.Fatalf("Error generating custom types: %v", err)
	}
	err = domain_gen.NewDomainGenerator(g.Domain).Generate()
	if err != nil {
		logrus.Fatalf("Error generating models: %v", err)
	}
	err = error_with_codes_gen.GenerateErrorWithCodes(g.ErrorWithCodes)
	if err != nil {
		logrus.Fatalf("Error generating error with codes: %v", err)
	}
	err = result_gen.GenerateResult(g.Result)
	if err != nil {
		logrus.Fatalf("Error generating result: %v", err)
	}
	err = config_gen.GenerateConfig(g.Config)
	if err != nil {
		logrus.Fatalf("Error generating config: %v", err)
	}
	err = env_gen.GenerateEnv(g.Env)
	if err != nil {
		logrus.Fatalf("Error generating env: %v", err)
	}
}

func formatingSettingsData(
	node map[string]interface{},
	currentPath string,
	folders *[]string,
	files *[]*custom_type_gen.CustomType,
	domains *[]*domain_gen.Domain,
	errorWithCodes *string,
	config *string,
	result *string,
	env *string,
) {
	for key, value := range node {
		newPath := filepath.Join(currentPath, key)

		switch v := value.(type) {
		case map[string]interface{}:
			// Обрабатываем вложенные директории
			*folders = append(*folders, newPath)
			if strings.Contains(newPath, "error_with_codes") {
				*errorWithCodes = newPath
			}
			if strings.Contains(newPath, "config") {
				*config = newPath
			}
			if strings.Contains(newPath, "result") {
				*result = newPath
			}
			if strings.Contains(newPath, "pkg/env") {
				*env = newPath
			}
			formatingSettingsData(v, newPath, folders, files, domains, errorWithCodes, config, result, env)

		case []interface{}:
			// Обрабатываем файлы
			*folders = append(*folders, newPath)
			if strings.Contains(newPath, "common") {
				for _, item := range v {
					fileData := item.(map[string]interface{})
					customType := custom_type_gen.NewCustomType(newPath, fileData["file_name"].(string), fileData["file_ext"].(string), fileData["file_type"].(string))
					*files = append(*files, customType)
				}
			}
			if strings.Contains(newPath, "model") || strings.Contains(newPath, "value_object") {
				for _, item := range v {
					fileData := item.(map[string]interface{})
					fields := make([]*domain_gen.Field, 0)
					for _, field := range fileData["fields"].([]interface{}) {
						fieldData := field.(map[string]interface{})
						fields = append(fields, &domain_gen.Field{
							Name: fieldData["name"].(string),
							Type: fieldData["type"].(string),
							Tag:  fieldData["tag"].(string),
						})
					}

					domain := domain_gen.NewDomain(newPath, fileData["file_name"].(string), fileData["file_ext"].(string), fields)
					*domains = append(*domains, domain)
				}
			}
		}
	}
}
