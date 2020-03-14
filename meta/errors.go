package meta

import (
	"errors"
)

// error(s)
var (
	ErrRepoError          = errors.New("Repo Error")
	ErrRepoRecordNotFound = errors.New("Record Not Found")

	ErrInvalidCode = errors.New("code invalid")

	ErrInvalidJWT = errors.New("jwt invalid")
	ErrSignJWT    = errors.New("jwt sign error")

	ErrWrongUsername = errors.New("Wrong Username/Phone/Email")
	ErrUserHasExist  = errors.New("Username/Phone/Email Has Exist")
	ErrCredential    = errors.New("Credential Error")

	ErrUnsupportedAuthType   = errors.New("Unsupported authorization method")
	ErrUnsupportedVerifyType = errors.New("Unsupported verify method")
	ErrUnsupportedActionType = errors.New("Unsupported action type")

	ErrAuthenticationFailed = errors.New("Authentication Failed")
	ErrAccessDenied         = errors.New("Access Denied")

	ErrAppHasExist   = errors.New("AppName Has Exist")
	ErrOrgHasExist   = errors.New("Org Has Exist")
	ErrOrgNotExist   = errors.New("Org Not Exist")
	ErrGroupHasExist = errors.New("Group Has Exist")
)
