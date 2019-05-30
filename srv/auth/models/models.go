package models

import (
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// BeforeCreateHook for gorm
type BeforeCreateHook struct{}

// BeforeCreate for gorm
func (*BeforeCreateHook) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("Id", uuid.String())
}

// UserInfo Model
type UserInfo struct {
	*common.Model
	*proto.UserInfo
	*BeforeCreateHook
}

// BeforeCreate for gorm
//func (*UserInfo) BeforeCreate(scope *gorm.Scope) error {
//	uuid := uuid.NewV4()
//	return scope.SetColumn("Id", uuid.String())
//}

// RoleInfo Model
type RoleInfo struct {
	*common.Model
	*proto.RoleInfo
	*BeforeCreateHook
}

// BeforeCreate for gorm
//func (*RoleInfo) BeforeCreate(scope *gorm.Scope) error {
//	uuid := uuid.NewV4()
//	return scope.SetColumn("Id", uuid.String())
//}

// PermissionInfo Model
type PermissionInfo struct {
	*common.Model
	*proto.PermissionInfo
	*BeforeCreateHook
}

// BeforeCreate for gorm
//func (*PermissionInfo) BeforeCreate(scope *gorm.Scope) error {
//	uuid := uuid.NewV4()
//	return scope.SetColumn("Id", uuid.String())
//}
