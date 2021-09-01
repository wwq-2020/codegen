package main

const (
	confGoTpl = "package conf\n import (\n\"stash.weimob.com/devops/go_common/rpc\"\n\"stash.weimob.com/devops/go_common/tracing\"\n\"stash.weimob.com/devops/go_common/confx\"\n\"stash.weimob.com/devops/go_common/gormx\"\n)\n // Conf Conf\ntype Conf struct {\nServer *rpc.ServerConf `toml:\"server\"`\nTracing *tracing.Conf `toml:\"tracing\"`\nDB *gormx.Conf `toml:\"db\"`\n} \n//MustParse MustParse\nfunc MustParse(cfgPath string) *Conf {\nconf:=&Conf{}\nconfx.MustParseFile(cfgPath,conf)\nreturn conf\n}\n"
)
