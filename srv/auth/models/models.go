package models

import (
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/common"
)

// User Model (for UserAuth)
type User struct {
	*common.Model
	*proto.UserAuth
	Roles []*Role `gorm:"many2many:user_roles;"`
}

// Role Model
type Role struct {
	*common.Model
	*proto.RoleInfo
	User        []*User       `gorm:"many2many:user_roles;"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

// Permission Model
type Permission struct {
	*common.Model
	*proto.PermissionInfo
	Roles []*Role `gorm:"many2many:role_permissions;"`
}

//// UserRole Model
//type UserRole struct {
//	*common.Model
//	ID string `gorm:"primary_key"`
//
//	UserID string `gorm:"foreignkey:CategoryId"`
//	RoleID string
//}
//
//// RolePermission Model
//type RolePermission struct {
//	*common.Model
//	ID string `gorm:"primary_key"`
//
//	RoleID       string
//	PermissionID string
//}
