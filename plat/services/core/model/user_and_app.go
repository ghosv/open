package model

import (
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User Model
type User struct {
	meta.Model

	Name  string `gorm:"not null;unique"`
	Phone string `gorm:"not null;unique"`
	Email string `gorm:"not null;unique"`
	Pass  string `gorm:"not null"`

	Nick     string `gorm:"not null"`
	Avatar   string `gorm:"not null"`
	Motto    string `gorm:"not null"`
	Homepage string `gorm:"not null"`

	MyApps []App `gorm:"ForeignKey:OwnerID"`   // 创建的应用
	Apps   []App `gorm:"many2many:app_users;"` // 桌面应用
}

// TokenPayload of User
func (u *User) TokenPayload() *pb.TokenPayload {
	return &pb.TokenPayload{
		Base: &pb.Identity{
			UUID:  u.ID,
			Name:  u.Name,
			Phone: u.Phone,
			Email: u.Email,
		},
		Info: &pb.UserInfo{
			UUID:     u.ID,
			Nick:     u.Nick,
			Avatar:   u.Avatar,
			Motto:    u.Motto,
			Homepage: u.Homepage,
		},
	}
}

// App Model
type App struct {
	meta.Model        // ID is app_id
	Key        string `gorm:"not null;unique"` // app_key (公匙/账号)
	Secret     string `gorm:"not null"`        // app_secret (私钥/密码)

	Name  string `gorm:"not null;unique"`
	Icon  string `gorm:"not null"`
	Intro string `gorm:"not null"`
	URL   string `gorm:"not null"` // <url>?token=<jwt>

	OwnerID    string        `gorm:"not null"`                   // "" => core TODO: change
	Managers   []User        `gorm:"many2many:app_managers;"`    // 管理者 TODO: change
	Developers []User        `gorm:"many2many:app_developers;"`  // 开发者 TODO: change
	Users      []User        `gorm:"many2many:app_users;"`       // 用户
	Scopes     []AccessScope `gorm:"ForeignKey:AppID"`           // 可申请范围
	AccessList []AccessScope `gorm:"many2many:app_access_list;"` // 申请范围
}

// BeforeCreate for gorm
func (*App) BeforeCreate(scope *gorm.Scope) error {
	if e := scope.SetColumn("ID", uuid.NewV4().String()); e != nil {
		return e
	}
	if e := scope.SetColumn("Key", uuid.NewV4().String()); e != nil {
		return e
	}
	if e := scope.SetColumn("Secret", uuid.NewV4().String()); e != nil {
		return e
	}
	return nil
}

// AccessScope Model
type AccessScope struct {
	gorm.Model
	AppID string

	// AppID:<scope>:<access> e.g. self:group:r
	Name   string `gorm:"not null;unique"`
	Detail string `gorm:"not null"`

	// Apps []AccessScope `gorm:"many2many:app_access_list;"`
}
