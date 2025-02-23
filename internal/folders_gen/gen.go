package folders_gen

import (
	"os"

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

func (fg *FolderGenerator) Generate() error {
	for _, folder := range fg.Folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			logrus.WithFields(logrus.Fields{
				"folder": folder,
				"error":  err,
			}).Error("Failed to create folder")
			return err
		}
	}

	return nil
}
