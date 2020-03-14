package utils

import (
	"github.com/ghosv/open/meta"
	"github.com/jinzhu/gorm"
)

// // RepoErrorFilter Utils - avoid importing gorm
// func RepoErrorFilter(e error) error {
// 	if e.Error() == "record not found" {
// 		return meta.ErrRepoRecordNotFound
// 	}
// 	return meta.ErrRepoError
// }

// RepoErrorFilter Utils - gate will not import this utils
func RepoErrorFilter(result *gorm.DB) error {
	if result.RecordNotFound() {
		return meta.ErrRepoRecordNotFound
	}
	return result.Error
}
