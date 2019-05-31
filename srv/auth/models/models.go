package models

import (
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/common"
)

// UserInfo Model
type UserInfo struct {
	*common.Model
	*proto.UserAuth
}

// RoleInfo Model
type RoleInfo struct {
	*common.Model
	*proto.RoleInfo
}

// PermissionInfo Model
type PermissionInfo struct {
	*common.Model
	*proto.PermissionInfo
}

// UserRole Model
type UserRole struct {
	*common.Model
	ID string `gorm:"primary_key"`

	UserID string
	RoleID string
}

// RolePermission Model
type RolePermission struct {
	*common.Model
	ID string `gorm:"primary_key"`

	RoleID       string
	PermissionID string
}
