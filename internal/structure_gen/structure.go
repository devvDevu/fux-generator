package structure_gen

import (
	"ca-generator/internal/gens/custom_type_gen"
	"ca-generator/internal/gens/domain_gen"
	"ca-generator/internal/gens/folders_gen"
	"ca-generator/internal/mocked_gens/config_gen"
	"ca-generator/internal/mocked_gens/env_gen"
	"ca-generator/internal/mocked_gens/error_with_codes_gen"
	"ca-generator/internal/mocked_gens/json_file_gen"
	"ca-generator/internal/mocked_gens/result_gen"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	wg := sync.WaitGroup{}
	// formating settings data
	formatingSettingsData(g.Settings["settings"].(map[string]interface{}), "", &g.Folders, &g.Files, &g.Domain, &g.ErrorWithCodes, &g.Config, &g.Result, &g.Env)

	// generating structure
	err := folders_gen.NewFolderGenerator(g.Folders).Generate()
	if err != nil {
		logrus.Fatalf("Error generating folders: %v", err)
	}

	// error catching goroutine
	errCh := make(chan error)
	go func() {
		for err := range errCh {
			logrus.Fatalf("Error generating: %v", err)
		}
	}()

	wg.Add(6)
	go func() {
		defer wg.Done()
		custom_type_gen.NewCustomTypeGenerator(g.Files).Generate(&wg, errCh)
	}()
	go func() {
		defer wg.Done()
		domain_gen.NewDomainGenerator(g.Domain).Generate(&wg, errCh)
	}()
	go func() {
		defer wg.Done()
		error_with_codes_gen.GenerateErrorWithCodes(&wg, errCh, g.ErrorWithCodes)
	}()
	go func() {
		defer wg.Done()
		result_gen.GenerateResult(&wg, errCh, g.Result)
	}()
	go func() {
		defer wg.Done()
		config_gen.GenerateConfig(&wg, errCh, g.Config)
	}()
	go func() {
		defer wg.Done()
		err := env_gen.GenerateEnv(g.Env)
		if err != nil {
			errCh <- err
		}
	}()
	wg.Wait()
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
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for key, value := range node {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newPath := filepath.Join(currentPath, key)
			mu.Lock()
			defer mu.Unlock()
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
		}()
	}
	wg.Wait()
}
