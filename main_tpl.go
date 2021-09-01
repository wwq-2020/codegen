package main

const (
	mainTpl = `
	package main

	import (
		"context"
		"flag"
		"stash.weimob.com/devops/go_common/rpc"
		"stash.weimob.com/devops/go_common/app"
		"stash.weimob.com/devops/go_common/log"
		"stash.weimob.com/devops/go_common/tracing"
		"{{$.ProjectPkg}}/pkg/conf"
		"{{.APIDocPkg}}"
	)
	var cfgPath = flag.String("config", "./conf/conf.toml", "-conf=./conf/conf.toml")
	
	func main() {
		flag.Parse()
		conf := conf.MustParse(*cfgPath)
		cleanup := tracing.MustInitGlobalTracer("{{$.ProjectName}}", conf.Tracing)
		defer cleanup()
		server := rpc.NewServer("{{$.ProjectName}}",conf.Server)
		{{range $idx,$each := .Services}}{{$each.Package}}Service:=Create{{$each.Package}}Service(conf)
		{{$.ProjectName}}.Register{{$each.Name}}RPCServer(server, {{$each.Package}}Service)
		{{end}}
		app.GoAsync(func(){
			if err := server.Start(); err != nil {
				log.WithError(err).
					Fatal("failed to Start")
			}
		})
		app.AddShutdownHook(func() {
			server.Stop(context.TODO())
		})
		app.Wait()
	}
	`
)
