package utils

import (
	"strings"

	metaPb "github.com/ghosv/open/meta/proto"
	pb "github.com/ghosv/open/plat/services/core/proto"
)

// ACL Util
type ACL struct {
	visitor *metaPb.Visitor
	appID   string
}

// NewACL visitor ask appID
func NewACL(token *pb.TokenPayload, appID string) *ACL {
	var scopes []string
	if token.Scopes == nil {
		scopes = []string{}
	} else {
		scopes = make([]string, 0, len(token.Scopes))
		for _, v := range token.Scopes {
			if !strings.HasPrefix(v.Name, appID+":") {
				continue
			}
			scopes = append(scopes, strings.Replace(v.Name, appID+":", "", 1))
		}
	}
	return &ACL{&metaPb.Visitor{
		UUID:   token.Base.UUID,
		AppID:  token.AppID,
		Scopes: scopes,
	}, appID}
}

// Visitor of ACL
func (a *ACL) Visitor() *metaPb.Visitor {
	return a.visitor
}

// Check Data <or>
func (a *ACL) Check(scopes ...string) bool {
	if a.visitor.AppID == "" || a.visitor.AppID == a.appID {
		// core's token or app's own token will pass
		return true
	}
	for _, v := range a.visitor.Scopes {
		for _, scope := range scopes {
			if v == scope {
				return true
			}
		}
	}
	return false
}

// Display uuid's Data <or>
func (a *ACL) Display(uuids ...string) bool {
	for _, v := range uuids {
		if a.visitor.UUID == v {
			return true
		}
	}
	return false
}
