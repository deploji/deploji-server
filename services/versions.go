package services

import (
	"errors"
	"fmt"
	"github.com/sotomskir/mastermind-server/models"
)

var GetVersions = func(appId uint) ([]models.Version, error) {
	app := models.GetApplication(appId)
	if app == nil {
		return nil, errors.New("not found")
	}
	var versions []models.Version
	if app.Repository.Type == "docker" {
		url := fmt.Sprintf("%s/v1/repositories/%s/tags", app.Repository.Url, app.RepositoryArtifact)
		var response []map[string]string
		err := GetJson(url, &response)
		if err != nil {
			return nil, err
		}
		for _, item := range response {
			versions = append(versions, models.Version{Name:item["name"]})
		}
	}
	return versions, nil
}
