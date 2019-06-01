package models

import (
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
)

// User Model (for UserAuth)
type User struct {
	*common.Model
	*proto.UserAuth
	Roles []*Role `gorm:"many2many:user_roles;"`
}

// Role Model
type Role struct {
	*gorm.Model
	*proto.RoleInfo
	User  []*User `gorm:"many2many:user_roles;"`
	Perms []*Perm `gorm:"many2many:role_perms;"`
}

// Perm Model
type Perm struct {
	*gorm.Model
	*proto.PermInfo
	Roles []*Role `gorm:"many2many:role_perms;"`
}
