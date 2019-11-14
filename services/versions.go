package services

import (
	"errors"
	"fmt"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
)

var GetVersions = func(appId uint) ([]dto.Version, error) {
	app := models.GetApplication(appId)
	if app == nil {
		return nil, errors.New("not found")
	}
	var versions []dto.Version
	if app.Repository.Type == "docker-v1" {
		url := fmt.Sprintf("%s/v1/repositories/%s/tags", app.Repository.Url, app.RepositoryArtifact)
		var response []map[string]string
		err := GetJson(url, &response)
		if err != nil {
			return nil, err
		}
		for _, item := range response {
			versions = append(versions, dto.Version{Name: item["name"]})
		}
		for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
			versions[i], versions[j] = versions[j], versions[i]
		}
	}
	if app.Repository.Type == "docker-v2" {
		url := fmt.Sprintf("%s/v2/%s/tags/list", app.Repository.Url, app.RepositoryArtifact)
		var response map[string]interface{}
		err := GetJson(url, &response)
		if err != nil {
			log.Printf("GetJson error: %s", err)
			return nil, err
		}

		for _, item := range response["tags"].([]interface{}) {
			versions = append(versions, dto.Version{Name: item.(string)})
		}
		for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
			versions[i], versions[j] = versions[j], versions[i]
		}
	}
	if app.Repository.Type == "nexus-v3" {
		continuationToken := ""
		hasMore := true
		versionsMap := make(map[string]dto.Version)
		regex := regexp.MustCompile(`-(\d{8})\.(\d{6})-(\d{1,4})`)

		for hasMore {
			url := fmt.Sprintf(
				"%s/service/rest/v1/search?repository=%s&group=%s&name=%s%s",
				app.Repository.Url,
				app.Repository.NexusName,
				app.RepositoryGroup,
				app.RepositoryArtifact,
				continuationToken)
			var response map[string]interface{}
			err := GetJson(url, &response)
			if err != nil {
				log.Printf("GetJson error: %s", err)
				return nil, err
			}
			hasMore = response["continuationToken"] != nil
			if hasMore {
				continuationToken = fmt.Sprintf("&continuationToken=%s", response["continuationToken"].(string))
			}

			for _, item := range response["items"].([]interface{}) {
				version := regex.ReplaceAllString(item.(map[string]interface{})["version"].(string), "-SNAPSHOT")
				snapshotData := regex.FindStringSubmatch(item.(map[string]interface{})["version"].(string))
				if len(snapshotData) == 4 {
					date, err := time.Parse("20060102", snapshotData[1])
					time, err2 := time.Parse("150405", snapshotData[2])
					if err == nil && err2 == nil {
						version = version + " | " + date.Format("2006-01-02") + " " + time.Format("15:04:05") + " | Build number: #" + snapshotData[3]
					}
				}

				versionsMap[version] = dto.Version{Name: version, SortKey: getSortKey(version)}
			}
		}
		for _, item := range versionsMap {
			versions = append(versions, item)
		}
		sort.Sort(dto.ByName(versions))
	}
	return versions, nil
}

func getSortKey(version string) string {
	regex := regexp.MustCompile(`\d{1,4}\.\d{1,4}\.\d{1,4}`)
	semver := string(regex.Find([]byte(version)))
	parts := strings.Split(semver, ".")
	if len(parts) < 3 {
		return version
	}
	return fmt.Sprintf("%03s.%03s.%03s", parts[0], parts[1], parts[2])
}
