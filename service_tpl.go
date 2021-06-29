package main

const (
	serviceTpl = `
	package service
	
	import (
		"context"
		"github.com/google/wire"
		"github.com/wwq-2020/go.common/rpc"
		"github.com/wwq-2020/go.common/app"
		"github.com/wwq-2020/go.common/log"
		"{{.ProjectPkg}}/pkg/{{.Package}}/repo"
		"{{.ProjectPkg}}/pkg/conf"
		"{{.APIDocPkg}}"
	)
	// Interface Interface
	type Interface interface {
		{{range $idx,$each := .APIs}}{{$each.Name}}(ctx context.Context, req *{{$.ProjectName}}.{{$each.Req}}) (*{{$.ProjectName}}.{{$each.Resp}}, error)
		{{end}}
	}

	type service struct {
		repo repo.Interface
	}

	// SuperSet SuperSet
	var SuperSet = wire.NewSet(MustNew, repo.MustNew, conf.MustNew)

	// MustNew MustNew
	func MustNew(repo repo.Interface) Interface{
		return &service{
			repo:repo,
		}
	}
	{{range $idx,$each := .APIs}}
	func(s *service) {{$each.Name}}(ctx context.Context, req *{{$.ProjectName}}.{{$each.Req}}) (*{{$.ProjectName}}.{{$each.Resp}}, error){
		return *{{$.ProjectName}}.{{$each.Resp}}{},nil
	}
	{{end}}
	`
)
