package main

const (
	mainTpl = `
	package main

	import (
		"context"
		"github.com/wwq-2020/go.common/rpc"
		"github.com/wwq-2020/go.common/app"
		"github.com/wwq-2020/go.common/log"
		"{{.APIDocPkg}}"
	)
	
	func main() {
		server := rpc.NewServer()
		{{range $idx,$each := .Services}}{{$each.Package}}Service:=Create{{$each.Package}}Service()
		{{$.ProjectName}}.Register{{$each.Name}}ToRPCServer(server, {{$each.Package}}Service){{end}}
		app.GoAsync(func(){
			if err := server.ListenAndServe(); err != nil {
				log.WithError(err).
					Fatal("failed to ListenAndServe")
			}
		})
		app.AddShutdownHook(func() {
			server.Stop(context.TODO())
		})
		app.Wait()
	}
	`
)
