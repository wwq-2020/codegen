package main

const (
	wireTpl = `
	//+build wireinject

	package main
	
	import (
		"github.com/google/wire"
		"github.com/wwq-2020/go.common/rpc"
		"github.com/wwq-2020/go.common/app"
		"github.com/wwq-2020/go.common/log"
		{{range $idx,$each := .Services}}"{{$.ProjectPkg}}/pkg/{{$each.Package}}"{{end}}
	)
	
	{{range $idx,$each := .Services}}
	// Create{{$each.Package}}Service Create{{$each.Package}}Service
	func Create{{$each.Package}}Service() {
		panic(wire.Build({{$each.Package}}.MustNew))
	}
	{{end}}
	`
)
