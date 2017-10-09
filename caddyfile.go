package caddyconsul

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/hashicorp/consul/api"
)

var consulGenerator *caddyfile

type caddyfile struct {
	contents    string
	lastKV      uint64
	lastService uint64
	domains     map[string]*domain
	services    map[string][]*api.CatalogService
}

func (s *caddyfile) Body() []byte {
	fmt.Println("Generated config:")
	fmt.Println(s.contents)
	return []byte(s.contents)
}

func (s *caddyfile) Path() string {
	return ""
}

func (s *caddyfile) ServerType() string {
	return "http"
}

func (s *caddyfile) StartWatching() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			s.WatchServices(true)
		}
	}()
}

func (s *caddyfile) buildConfig() {
	t := template.New("Caddyfile.tmpl")
	var err error
	t, err = t.ParseFiles(caddyfilePath)
	if err != nil {
		s.contents = ""
		return
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, s.services); err != nil {
		return
	}
	s.contents = tpl.String()
}
