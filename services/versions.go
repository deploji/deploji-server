package services

import (
	"errors"
	"fmt"
	"github.com/sotomskir/mastermind-server/dto"
	"github.com/sotomskir/mastermind-server/models"
	"log"
)

var GetVersions = func(appId uint) ([]dto.Version, error) {
	app := models.GetApplication(appId)
	if app == nil {
		return nil, errors.New("not found")
	}
	var versions []dto.Version
	if app.Repository.Type == "docker" {
		url := fmt.Sprintf("%s/v2/%s/tags/list", app.Repository.Url, app.RepositoryArtifact)
		var response map[string]interface{}
		err := GetJson(url, &response)
		if err != nil {
			log.Printf("GetJson error: %s", err)
			return nil, err
		}

		for _, item := range response["tags"].([]interface{}) {
			versions = append(versions, dto.Version{Name:item.(string)})
		}
	}
	return versions, nil
}
