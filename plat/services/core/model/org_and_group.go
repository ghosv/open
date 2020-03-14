package model

import (
	"github.com/ghosv/open/meta"
)

// Org Model
type Org struct {
	meta.Model
	FatherID string

	Name   string `gorm:"not null"`
	Icon   string `gorm:"not null"`
	Detail string `gorm:"not null"`

	MasterID string // UUID
	Master   User
	Users    []User `gorm:"many2many:org_users;"`
}

// Group Model
type Group struct {
	meta.Model

	Name   string `gorm:"not null"`
	Icon   string `gorm:"not null"`
	Detail string `gorm:"not null"`

	MasterID string // UUID
	Master   User
	Users    []User `gorm:"many2many:group_users;"`
}

// TODO: User 的 Org/Group 对 App 授权可见
