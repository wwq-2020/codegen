package main

const (
	confTpl = "package conf\n// Conf Conf\ntype Conf struct {\nServer *Server `toml:\"server\"`\n}\n\n// Server Server\ntype Server struct {\nAddr string\n}\n"
)
