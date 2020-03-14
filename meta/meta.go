package meta

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// meta
const (
	SystemName = "GhostOpenPlatform"

	// Version = "1.0.0"
	Version = "latest"

	// === Client ===
	ClientGateway = "open.gateway"

	// === Platform Services ===
	SrvCore = "open.srv.core"
	// SrvManage = "open.srv.manage"
	SrvSelf   = "open.srv.self"
	SrvNotify = "open.srv.notify"

	// === Apps ===
	AppMrms = "open.app.mrms"
	// AppNote = "open.app.note"
)

// Model without ID
type Model struct {
	ID string `gorm:"primary_key"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate for meta.Model
func (*Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4().String())
}
