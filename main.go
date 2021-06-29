package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wwq-2020/go.common/log"
)

var (
	serviceReg = regexp.MustCompile("func Register(.*)ToRPCServer")
	apiReg     = regexp.MustCompile("\t(.*)\\(context.Context, \\*(.*)\\) \\(\\*(.*), error\\)")
)

type mainData struct {
	ProjectName string
	APIDocPkg   string
	ProjectPkg  string
	Services    []*service
}

type service struct {
	APIDocPkg   string
	ProjectPkg  string
	ProjectName string
	Name        string
	Package     string
	APIs        []*api
}

type api struct {
	Name string
	Req  string
	Resp string
}

type conf struct {
	APIDocGit        string `json:"apidoc_git"`
	Dir              string `json:"dir"`
	APIDocPkg        string `json:"apidoc_pkg"`
	ProjectPkgPrefix string `json:"project_pkg_prefix"`
}

func main() {
	if len(os.Args) <= 1 {
		log.Info("missing project name")
		return
	}
	projectName := os.Args[1]

	home, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).
			Fatal("failed to UserHomeDir")
	}
	_confFile := path.Join(home, ".codegen")
	confData, err := ioutil.ReadFile(_confFile)
	if err != nil {
		log.WithError(err).
			Fatal("failed to ReadFile")
	}
	conf := &conf{}
	if err := json.Unmarshal(confData, conf); err != nil {
		log.WithError(err).
			Fatal("failed to Unmarshal")
	}
	flag.Parse()

	tempDir := os.TempDir()
	apidocDir := path.Join(tempDir, "apidoc")
	if err := os.RemoveAll(apidocDir); err != nil {
		log.WithError(err).
			Fatal("failed to Run")
	}
	cmd := exec.Command("git", "clone", conf.APIDocGit, apidocDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.WithError(err).
			Fatal("failed to Run")
	}
	root := path.Join(apidocDir, conf.Dir, projectName)
	var services []*service
	if err := filepath.WalkDir(root, func(curPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(curPath)
		if err != nil {
			return err
		}
		serviceData := serviceReg.FindStringSubmatch(string(data))
		if len(serviceData) != 2 {
			log.WithField("path", curPath).
				Info("find no service")
			return nil
		}
		service := &service{
			APIDocPkg:   path.Join(conf.APIDocPkg, conf.Dir, projectName),
			ProjectPkg:  path.Join(conf.ProjectPkgPrefix, projectName),
			ProjectName: projectName,
		}
		service.Name = serviceData[1]
		service.Package = strings.ToLower(serviceData[1])
		apiData := apiReg.FindAllStringSubmatch(string(data), -1)
		if len(apiData) == 0 {
			log.WithField("path", curPath).
				Info("find no api")
			return nil
		}
		for _, each := range apiData {
			service.APIs = append(service.APIs, &api{
				Name: each[1],
				Req:  each[2],
				Resp: each[3],
			})
		}
		services = append(services, service)
		return nil
	}); err != nil {
		log.WithError(err).
			Fatal("failed to Walk")
	}
	cmdDir := path.Join(projectName, "cmd")
	if err = os.MkdirAll(cmdDir, 0744); err != nil {
		log.WithError(err).
			Fatal("failed to MkdirAll")
	}
	buf := bytes.NewBuffer(nil)
	mainData := &mainData{
		APIDocPkg:   path.Join(conf.APIDocPkg, conf.Dir, projectName),
		ProjectPkg:  path.Join(conf.ProjectPkgPrefix, projectName),
		ProjectName: projectName,
		Services:    services,
	}

	t, err := template.New("main_tpl").Funcs(template.FuncMap{}).Parse(mainTpl)
	if err != nil {
		log.WithError(err).
			Fatal("failed to Parse")
	}
	if err := t.Execute(buf, mainData); err != nil {
		log.WithError(err).
			Fatal("failed to Execute")
	}
	mainFile := path.Join(cmdDir, "main.go")
	if err := ioutil.WriteFile(mainFile, buf.Bytes(), 0644); err != nil {
		log.WithError(err).
			Fatal("failed to WriteFile")
	}
	buf.Reset()

	t, err = template.New("wire_tpl").Funcs(template.FuncMap{}).Parse(wireTpl)
	if err != nil {
		log.WithError(err).
			Fatal("failed to Parse")
	}
	if err := t.Execute(buf, mainData); err != nil {
		log.WithError(err).
			Fatal("failed to Execute")
	}
	wireFile := path.Join(cmdDir, "wire.go")
	if err := ioutil.WriteFile(wireFile, buf.Bytes(), 0644); err != nil {
		log.WithError(err).
			Fatal("failed to WriteFile")
	}
	buf.Reset()
	t, err = template.New("conf_tpl").Funcs(template.FuncMap{}).Parse(confTpl)
	if err != nil {
		log.WithError(err).
			Fatal("failed to Parse")
	}
	if err := t.Execute(buf, mainData); err != nil {
		log.WithError(err).
			Fatal("failed to Execute")
	}
	confDir := path.Join(projectName, "pkg", "conf")
	if err := os.MkdirAll(confDir, 0744); err != nil {
		log.WithError(err).
			Fatal("failed to MkdirAll")
	}
	confFile := path.Join(confDir, "conf.go")

	if err := ioutil.WriteFile(confFile, buf.Bytes(), 0644); err != nil {
		log.WithError(err).
			Fatal("failed to WriteFile")
	}

	buf.Reset()
	for _, service := range mainData.Services {
		serviceDir := path.Join(projectName, "pkg", service.Package, "service")
		if err := os.MkdirAll(serviceDir, 0744); err != nil {
			log.WithError(err).
				Fatal("failed to MkdirAll")
		}
		buf.Reset()
		t, err = template.New("service_tpl").Funcs(template.FuncMap{}).Parse(serviceTpl)
		if err != nil {
			log.WithError(err).
				Fatal("failed to Parse")
		}
		if err := t.Execute(buf, service); err != nil {
			log.WithError(err).
				Fatal("failed to Execute")
		}

		serviceFile := path.Join(serviceDir, "service.go")

		if err := ioutil.WriteFile(serviceFile, buf.Bytes(), 0644); err != nil {
			log.WithError(err).
				Fatal("failed to WriteFile")
		}
		repoDir := path.Join(projectName, "pkg", service.Package, "repo")
		if err := os.MkdirAll(repoDir, 0744); err != nil {
			log.WithError(err).
				Fatal("failed to MkdirAll")
		}
		buf.Reset()
		t, err = template.New("repo_tpl").Funcs(template.FuncMap{}).Parse(repoTpl)
		if err != nil {
			log.WithError(err).
				Fatal("failed to Parse")
		}
		if err := t.Execute(buf, mainData); err != nil {
			log.WithError(err).
				Fatal("failed to Execute")
		}

		repoFile := path.Join(repoDir, "repo.go")

		if err := ioutil.WriteFile(repoFile, buf.Bytes(), 0644); err != nil {
			log.WithError(err).
				Fatal("failed to WriteFile")
		}
		modelDir := path.Join(projectName, "pkg", service.Package, "model")
		if err := os.MkdirAll(modelDir, 0744); err != nil {
			log.WithError(err).
				Fatal("failed to MkdirAll")
		}
	}
	if err := os.Chdir(cmdDir); err != nil {
		log.WithError(err).
			Fatal("failed to Chdir")
	}
	cmd = exec.Command("wire")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.WithError(err).
			Fatal("failed to wire")
	}

	cmd = exec.Command("go", "mod", "init", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmdDir, "------", err)

}
