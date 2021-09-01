package main

const (
	repoTpl = `
	package repo

import (
	"{{.ProjectPkg}}/pkg/conf"
	"gorm.io/gorm"
	"stash.weimob.com/devops/go_common/gormx"
)

// Interface Interface
type Interface interface {
}

type repo struct{
	db *gorm.DB
}

// MustNew MustNew
func MustNew(conf *conf.Conf) Interface {
	db := gormx.MustOpen(conf.DB)
	return &repo{
		db:db,
	}
}
`
)
