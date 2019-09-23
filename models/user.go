package models

import (
	"github.com/jinzhu/gorm"
	"log"
)

type Password string
type UserType string

const (
	UserTypeAdmin   UserType = "admin"
	UserTypeAuditor UserType = "auditor"
	UserTypeRegular UserType = "regular"
)

type User struct {
	gorm.Model
	Name     string
	Surname  string
	Username string `gorm:"unique_index"`
	Email    string
	Password Password
	IsActive bool
	Type     UserType
}

// Marshaler ignores the field value completely.
func (Password) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func GetUsers() []*User {
	users := make([]*User, 0)
	err := GetDB().Find(&users).Error
	if err != nil {
		return nil
	}
	return users
}

func GetUser(id uint) *User {
	var user User
	err := GetDB().First(&user, id).Error
	if err != nil {
		return nil
	}
	return &user
}

func GetUserByUsername(username string) *User {
	var user User
	err := GetDB().Where("username=?", username).Find(&user).Error
	if err != nil {
		log.Printf("User not found: %s", err)
		return nil
	}
	return &user
}

func SaveUser(user *User) error {
	if GetDB().NewRecord(user) {
		err := GetDB().Create(user).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(user).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteUser(user *User) error {
	err := GetDB().Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}
