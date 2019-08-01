package services

import (
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

type TestType struct {
	Id   uint `json:"id"`
	Name string `json:"name"`
}

func TestGetJson(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	response := `{"id": "11", "name": "mastermind"}`
	httpmock.RegisterResponder("GET", "https://example.com/test.json", httpmock.NewStringResponder(200, response))
	var json map[string]string
	err := GetJson("https://example.com/test.json", &json)
	if err != nil {
		t.Errorf("TestGetJson: Error thrown: %s", err)
	}
	if json["id"] != "11" {
		t.Errorf("TestGetJson: expected: 11, got: %s", json["id"])
	}
	if json["name"] != "mastermind" {
		t.Errorf("TestGetJson: expected: mastermind, got: %s", json["name"])
	}
}
