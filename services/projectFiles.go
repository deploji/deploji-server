package services

import (
	"errors"
	"fmt"
	"github.com/sotomskir/mastermind-server/models"
	"os"
	"path/filepath"
)

var GetProjectFiles = func(projectId uint) ([]models.ProjectFile, error) {
	project := models.GetProject(projectId)
	if project == nil {
		return nil, errors.New("not found")
	}
	var projectFiles []models.ProjectFile
	root := fmt.Sprintf("./storage/repositories/%s", project.Name)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			projectFiles = append(projectFiles, models.ProjectFile(path[len(root)-1:]))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return projectFiles, nil
}
