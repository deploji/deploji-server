package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"strconv"
	"strings"
)

var E *casbin.SyncedEnforcer

type PermissionType string
type ActionType string

const (
	PermissionTypeInventory    PermissionType = "inventories"
	PermissionTypeSshKey       PermissionType = "ssh-keys"
	PermissionTypeTemplate     PermissionType = "templates"
	PermissionTypeApplications PermissionType = "applications"
	PermissionTypeTeams        PermissionType = "teams"
	PermissionTypeUsers        PermissionType = "users"
	PermissionTypeJobs         PermissionType = "jobs"
	PermissionTypeRepositories PermissionType = "repositories"
)

const (
	ActionTypeRead  ActionType = "GET"
	ActionTypeWrite ActionType = "POST"
	ActionTypeAdmin ActionType = "ADMIN"
)

func GetPermissionsForTeam(teamId string) []dto.Permission {
	permissions := make([]dto.Permission, 0)
	for _, permission := range E.GetPermissionsForUser(fmt.Sprintf("%s/%s", PermissionTypeTeams, teamId)) {
		objectType, id := splitPermission(permission[1])
		objectName := getObjectName(objectType, id)
		permissions = append(permissions, dto.Permission{
			Name:       objectName,
			ObjectType: string(objectType),
			ObjectID:   id,
			Role:       permission[2],
		})
	}
	return permissions
}

func getObjectName(permissionType PermissionType, id uint) string {
	var name string
	switch permissionType {
	case PermissionTypeInventory:
		inventory := models.GetInventory(id)
		if inventory != nil {
			name = inventory.Name
		}
	case PermissionTypeTemplate:
		template := models.GetTemplate(id)
		if template != nil {
			name = template.Name
		}
	case PermissionTypeApplications:
		application := models.GetApplication(id)
		if application != nil {
			name = application.Name
		}
	case PermissionTypeSshKey:
		key := models.GetSshKey(id)
		if key != nil {
			name = key.Title
		}
	default:
		name = fmt.Sprintf("%d", id)
	}
	return fmt.Sprintf("%s: %s", permissionType, name)
}

func splitPermission(perm string) (PermissionType, uint) {
	split := strings.Split(perm, "/")
	permType := PermissionType(split[0])
	id, _ := strconv.ParseUint(split[1], 10, 16)
	return permType, uint(id)
}

func GetUsersForRole(teamId string) []models.User {
	casbinUsers, _ := E.GetUsersForRole(fmt.Sprintf("%s/%s", PermissionTypeTeams, teamId))
	users := make([]models.User, 0)
	for _, username := range casbinUsers {
		id, _ := strconv.ParseUint(strings.TrimLeft(username, fmt.Sprintf("%s/", PermissionTypeUsers)), 10, 16)
		user := models.GetUser(uint(id))
		if user != nil {
			users = append(users, *user)
		}
	}
	return users
}

func AddUserToTeam(teamId string, userId uint) error {
	_, err := E.AddGroupingPolicy(fmt.Sprintf("%s/%d", PermissionTypeUsers, userId), fmt.Sprintf("%s/%s", PermissionTypeTeams, teamId))
	return err
}

func RemoveUserFromTeam(teamId string, userId string) error {
	_, err := E.RemoveGroupingPolicy(fmt.Sprintf("%s/%s", PermissionTypeUsers, userId), fmt.Sprintf("%s/%s", PermissionTypeTeams, teamId))
	return err
}

func AddPermissionToTeam(teamId string, objectType string, objectId uint, role string) error {
	_, err := E.AddPolicy(fmt.Sprintf("%s/%s", PermissionTypeTeams, teamId), fmt.Sprintf("%s/%d", objectType, objectId), role)
	return err
}

func RemovePermissionFromTeam(teamId string, objectType string, objectId uint, role string) error {
	_, err := E.RemovePolicy(fmt.Sprintf("%s/%s", PermissionTypeTeams, teamId), fmt.Sprintf("%s/%d", objectType, objectId), role)
	return err
}

func GetImplicitPermissionsForUser(id uint) ([][]string, error) {
	return E.GetImplicitPermissionsForUser(fmt.Sprintf("%s/%d", PermissionTypeUsers, id))
}

func Enforce(user dto.JWTClaims, permType PermissionType, id uint, actionType ActionType) bool {
	if user.Type == models.UserTypeAdmin {
		return true
	}
	if user.Type == models.UserTypeAuditor && actionType == ActionTypeRead {
		return true
	}
	if permType == PermissionTypeRepositories ||
		permType == PermissionTypeJobs {
		return true
	}
	sub := fmt.Sprintf("%s/%d", PermissionTypeUsers, user.UserID)
	obj := fmt.Sprintf("%s/%d", permType, id)
	isAllowed, err := E.Enforce(sub, obj, string(actionType))
	if err != nil {
		return false
	}
	if !isAllowed && actionType == ActionTypeRead {
		isAllowed, err = E.Enforce(sub, obj, string(ActionTypeAdmin))
		if err != nil {
			return false
		}
	}
	return isAllowed
}

func FilterTemplates(templates []*models.Template, user dto.JWTClaims) []*models.Template {
	result := make([]*models.Template, 0)
	for _, template := range templates {
		if Enforce(user, PermissionTypeTemplate, template.ID, ActionTypeRead) {
			result = append(result, template)
		}
	}
	return result
}

func FilterSshKeys(keys []*models.SshKey, user dto.JWTClaims) []*models.SshKey {
	result := make([]*models.SshKey, 0)
	for _, key := range keys {
		if Enforce(user, PermissionTypeSshKey, key.ID, ActionTypeRead) {
			result = append(result, key)
		}
	}
	return result
}

func FilterInventories(inventories []*models.Inventory, user dto.JWTClaims) []*models.Inventory {
	result := make([]*models.Inventory, 0)
	for _, inventory := range inventories {
		if Enforce(user, PermissionTypeInventory, inventory.ID, ActionTypeRead) {
			result = append(result, inventory)
		}
	}
	return result
}

func FilterApplications(applications []*models.Application, user dto.JWTClaims) []*models.Application {
	result := make([]*models.Application, 0)
	for _, app := range applications {
		if Enforce(user, PermissionTypeApplications, app.ID, ActionTypeRead) {
			result = append(result, app)
		}
	}
	return result
}
