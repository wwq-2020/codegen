package main

const (
	confTpl = "[server]\naddr= \"127.0.0.1:8080\"\n[tracing]endpoint=\"http://127.0.0.1:14268/api/traces?format=jaeger.thrift\"\nsample_rate=1.0\n"
)
