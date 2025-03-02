package folders_gen

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type FolderGenerator struct {
	Folders []string
}

func NewFolderGenerator(folders []string) *FolderGenerator {
	return &FolderGenerator{
		Folders: folders,
	}
}

func (fg *FolderGenerator) Generate(wg *sync.WaitGroup) error {
	for _, folder := range fg.Folders {
		wg.Add(1)
		go func(folder string) {
			defer wg.Done()
			if err := os.MkdirAll(folder, 0755); err != nil {
				logrus.WithFields(logrus.Fields{
					"folder": folder,
					"error":  err,
				}).Error("Failed to create folder")
			}
		}(folder)
	}
	return nil
}
