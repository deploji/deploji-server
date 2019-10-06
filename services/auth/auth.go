package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/utils"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var E *casbin.SyncedEnforcer

func getPermissionsForSubject(subjectType dto.SubjectType, subjectId uint) []dto.Permission {
	permissions := make([]dto.Permission, 0)
	for _, permission := range E.GetPermissionsForUser(fmt.Sprintf("%s/%d", subjectType, subjectId)) {
		objectType, id := splitPermission(permission[1])
		objectName := getObjectName(dto.ObjectType(objectType), id)
		subjectName := getSubjectName(subjectType, subjectId)
		permissions = append(permissions, dto.Permission{
			ObjectType:  dto.ObjectType(objectType),
			ObjectID:    id,
			ObjectName:  objectName,
			SubjectType: subjectType,
			SubjectID:   subjectId,
			SubjectName: subjectName,
			Action:      dto.ActionType(permission[2]),
		})
	}
	return permissions
}

func GetPermissions(filters []utils.Filter) []dto.Permission {
	permissions := make([]dto.Permission, 0)
	for _, permission := range E.GetPolicy() {
		subjectType, subjectId := splitPermission(permission[0])
		objectType, objectId := splitPermission(permission[1])
		if !matchFilters(filters, subjectType, subjectId, objectType, objectId) {
			continue
		}
		action := permission[2]
		objectName := getObjectName(dto.ObjectType(objectType), objectId)
		subjectName := getSubjectName(dto.SubjectType(subjectType), subjectId)
		permissions = append(permissions, dto.Permission{
			ObjectType:  dto.ObjectType(objectType),
			ObjectID:    objectId,
			ObjectName:  objectName,
			SubjectType: dto.SubjectType(subjectType),
			SubjectID:   subjectId,
			SubjectName: subjectName,
			Action:      dto.ActionType(action),
		})
	}
	return permissions
}

func matchFilters(filters []utils.Filter, subjectType string, subjectId uint, objectType string, objectId uint) bool {
	for _, filter := range filters {
		switch filter.Key {
		case "SubjectType":
			if subjectType != filter.Value {
				return false
			}
			break
		case "ObjectType":
			if objectType != filter.Value {
				return false
			}
		case "ObjectID":
			if fmt.Sprintf("%d", objectId) != filter.Value {
				return false
			}
		}
	}
	return true
}

func GetPermissionsForTeam(teamId uint) []dto.Permission {
	return getPermissionsForSubject(dto.SubjectTypeTeam, teamId)
}

func GetPermissionsForUser(userId uint) []dto.Permission {
	return getPermissionsForSubject(dto.SubjectTypeUser, userId)
}

func getObjectName(permissionType dto.ObjectType, id uint) string {
	var name string
	switch permissionType {
	case dto.ObjectTypeInventory:
		inventory := models.GetInventory(id)
		if inventory != nil {
			name = inventory.Name
		}
	case dto.ObjectTypeTemplate:
		template := models.GetTemplate(id)
		if template != nil {
			name = template.Name
		}
	case dto.ObjectTypeApplications:
		application := models.GetApplication(id)
		if application != nil {
			name = application.Name
		}
	case dto.ObjectTypeSshKey:
		key := models.GetSshKey(id)
		if key != nil {
			name = key.Title
		}
	default:
		name = fmt.Sprintf("%d", id)
	}
	return name
}

func getSubjectName(subjectType dto.SubjectType, id uint) string {
	var name string
	switch subjectType {
	case dto.SubjectTypeUser:
		user := models.GetUser(id)
		if user != nil {
			name = user.Username
		}
	case dto.SubjectTypeTeam:
		team := models.GetTeam(id)
		if team != nil {
			name = team.Name
		}
	default:
		name = fmt.Sprintf("%d", id)
	}
	return name
}

func splitPermission(perm string) (string, uint) {
	split := strings.Split(perm, "/")
	id, _ := strconv.ParseUint(split[1], 10, 16)
	return split[0], uint(id)
}

func GetUsersForTeam(teamId string) []models.User {
	casbinUsers, _ := E.GetUsersForRole(fmt.Sprintf("%s/%s", dto.SubjectTypeTeam, teamId))
	users := make([]models.User, 0)
	for _, username := range casbinUsers {
		id, _ := strconv.ParseUint(strings.TrimLeft(username, fmt.Sprintf("%s/", dto.SubjectTypeUser)), 10, 16)
		user := models.GetUser(uint(id))
		if user != nil {
			users = append(users, *user)
		}
	}
	return users
}

func AddUserToTeam(teamId string, userId uint) error {
	_, err := E.AddGroupingPolicy(
		fmt.Sprintf("%s/%d", dto.ObjectTypeUsers, userId),
		fmt.Sprintf("%s/%s", dto.ObjectTypeTeams, teamId))
	return err
}

func RemoveUserFromTeam(teamId string, userId string) error {
	_, err := E.RemoveGroupingPolicy(
		fmt.Sprintf("%s/%s", dto.ObjectTypeUsers, userId),
		fmt.Sprintf("%s/%s", dto.ObjectTypeTeams, teamId))
	return err
}

func AddPermission(permission dto.Permission) error {
	_, err := E.AddPolicy(
		fmt.Sprintf("%s/%d", permission.SubjectType, permission.SubjectID),
		fmt.Sprintf("%s/%d", permission.ObjectType, permission.ObjectID),
		string(permission.Action))
	return err
}

func RemovePermission(permission dto.Permission) error {
	_, err := E.RemovePolicy(
		fmt.Sprintf("%s/%d", permission.SubjectType, permission.SubjectID),
		fmt.Sprintf("%s/%d", permission.ObjectType, permission.ObjectID),
		string(permission.Action))
	return err
}

func AddOwnerPermissions(r *http.Request, object interface{}) {
	user := services.GetJWTClaims(r)
	if r.Method != http.MethodPost || user.Type == models.UserTypeAdmin {
		return
	}
	var objectType dto.ObjectType
	var objectID uint
	t := reflect.TypeOf(object)
	if t.ConvertibleTo(reflect.TypeOf(models.Inventory{})) {
		objectType = dto.ObjectTypeInventory
		objectID = object.(models.Inventory).ID
	}
	if t.ConvertibleTo(reflect.TypeOf(models.Application{})) {
		objectType = dto.ObjectTypeApplications
		objectID = object.(models.Application).ID
	}
	if t.ConvertibleTo(reflect.TypeOf(models.SshKey{})) {
		objectType = dto.ObjectTypeSshKey
		objectID = object.(models.SshKey).ID
	}
	if t.ConvertibleTo(reflect.TypeOf(models.Template{})) {
		objectType = dto.ObjectTypeTemplate
		objectID = object.(models.Template).ID
	}
	permission := dto.Permission{
		SubjectType: dto.SubjectTypeUser,
		SubjectID:   user.UserID,
		ObjectType:  objectType,
		ObjectID:    objectID,
		Action:      dto.ActionTypeAdmin,
	}
	_ = AddPermission(permission)
	permission.Action = dto.ActionTypeRead
	_ = AddPermission(permission)
	permission.Action = dto.ActionTypeWrite
	_ = AddPermission(permission)
}

func GetImplicitPermissionsForUser(id uint) ([][]string, error) {
	return E.GetImplicitPermissionsForUser(fmt.Sprintf("%s/%d", dto.ObjectTypeUsers, id))
}

func Enforce(user dto.JWTClaims, permType dto.ObjectType, id uint, actionType dto.ActionType) bool {
	if user.Type == models.UserTypeAdmin {
		return true
	}
	if user.Type == models.UserTypeAuditor && actionType == dto.ActionTypeRead {
		return true
	}
	if permType == dto.ObjectTypeRepositories ||
		permType == dto.ObjectTypeJobs {
		return true
	}
	sub := fmt.Sprintf("%s/%d", dto.ObjectTypeUsers, user.UserID)
	obj := fmt.Sprintf("%s/%d", permType, id)
	isAllowed, err := E.Enforce(sub, obj, string(actionType))
	if err != nil {
		return false
	}
	if !isAllowed && actionType == dto.ActionTypeRead {
		isAllowed, err = E.Enforce(sub, obj, string(dto.ActionTypeAdmin))
		if err != nil {
			return false
		}
	}
	return isAllowed
}

func InsertTemplatePermissions(object *models.Template, user dto.JWTClaims) {
	object.Read = Enforce(user, dto.ObjectTypeTemplate, object.ID, dto.ActionTypeRead)
	object.Write = Enforce(user, dto.ObjectTypeTemplate, object.ID, dto.ActionTypeWrite)
	object.Admin = Enforce(user, dto.ObjectTypeTemplate, object.ID, dto.ActionTypeAdmin)
}

func InsertApplicationPermissions(object *models.Application, user dto.JWTClaims) {
	object.Read = Enforce(user, dto.ObjectTypeApplications, object.ID, dto.ActionTypeRead)
	object.Write = Enforce(user, dto.ObjectTypeApplications, object.ID, dto.ActionTypeWrite)
	object.Admin = Enforce(user, dto.ObjectTypeApplications, object.ID, dto.ActionTypeAdmin)
}

func InsertSshKeyPermissions(object *models.SshKey, user dto.JWTClaims) {
	object.Read = Enforce(user, dto.ObjectTypeSshKey, object.ID, dto.ActionTypeRead)
	object.Write = Enforce(user, dto.ObjectTypeSshKey, object.ID, dto.ActionTypeWrite)
	object.Admin = Enforce(user, dto.ObjectTypeSshKey, object.ID, dto.ActionTypeAdmin)
}

func InsertInventoryPermissions(object *models.Inventory, user dto.JWTClaims) {
	object.Read = Enforce(user, dto.ObjectTypeInventory, object.ID, dto.ActionTypeRead)
	object.Write = Enforce(user, dto.ObjectTypeInventory, object.ID, dto.ActionTypeWrite)
	object.Admin = Enforce(user, dto.ObjectTypeInventory, object.ID, dto.ActionTypeAdmin)
}

func FilterTemplates(templates []*models.Template, user dto.JWTClaims) []*models.Template {
	result := make([]*models.Template, 0)
	for _, template := range templates {
		if Enforce(user, dto.ObjectTypeTemplate, template.ID, dto.ActionTypeRead) {
			InsertTemplatePermissions(template, user)
			result = append(result, template)
		}
	}
	return result
}

func FilterSshKeys(keys []*models.SshKey, user dto.JWTClaims) []*models.SshKey {
	result := make([]*models.SshKey, 0)
	for _, key := range keys {
		if Enforce(user, dto.ObjectTypeSshKey, key.ID, dto.ActionTypeRead) {
			InsertSshKeyPermissions(key, user)
			result = append(result, key)
		}
	}
	return result
}

func FilterInventories(inventories []*models.Inventory, user dto.JWTClaims) []*models.Inventory {
	result := make([]*models.Inventory, 0)
	for _, inventory := range inventories {
		if Enforce(user, dto.ObjectTypeInventory, inventory.ID, dto.ActionTypeRead) {
			InsertInventoryPermissions(inventory, user)
			result = append(result, inventory)
		}
	}
	return result
}

func FilterApplications(applications []*models.Application, user dto.JWTClaims) []*models.Application {
	result := make([]*models.Application, 0)
	for _, app := range applications {
		if Enforce(user, dto.ObjectTypeApplications, app.ID, dto.ActionTypeRead) {
			InsertApplicationPermissions(app, user)
			result = append(result, app)
		}
	}
	return result
}
