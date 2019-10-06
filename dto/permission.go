package dto

type SubjectType string
type ObjectType string
type ActionType string

const (
	SubjectTypeTeam SubjectType = "teams"
	SubjectTypeUser SubjectType = "users"
)

const (
	ObjectTypeInventory    ObjectType = "inventories"
	ObjectTypeSshKey       ObjectType = "ssh-keys"
	ObjectTypeTemplate     ObjectType = "templates"
	ObjectTypeApplications ObjectType = "applications"
	ObjectTypeTeams        ObjectType = "teams"
	ObjectTypeUsers        ObjectType = "users"
	ObjectTypeJobs         ObjectType = "jobs"
	ObjectTypeRepositories ObjectType = "repositories"
)

const (
	ActionTypeRead  ActionType = "read"
	ActionTypeWrite ActionType = "write"
	ActionTypeAdmin ActionType = "admin"
	ActionTypeUse   ActionType = "use"
)

type Permission struct {
	SubjectType SubjectType
	SubjectID   uint
	SubjectName string
	ObjectType  ObjectType
	ObjectID    uint
	ObjectName  string
	Action      ActionType
}
