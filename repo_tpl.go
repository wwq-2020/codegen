package main

const (
	repoTpl = `
	package repo

import "{{.ProjectPkg}}/pkg/conf"

// Interface Interface
type Interface interface {
}

type repo struct{}

// MustNew MustNew
func MustNew(conf *conf.Conf) Interface {
	return &repo{}
}
`
)
