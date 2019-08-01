package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/sotomskir/mastermind-server/dto"
)

var E *casbin.Enforcer

func GetGroups() []dto.Group {
	groups := make([]dto.Group, 0)
	for _, group := range E.GetAllSubjects() {
		groups = append(groups, dto.Group{Name: group})
	}
	return groups
}

func GetGroup(id string) {
	E.GetUsersForRole()
}
