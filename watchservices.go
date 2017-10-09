package caddyconsul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

type service struct {
	Name      string
	Instances []*api.CatalogService
}

func (s *caddyfile) WatchServices(reload bool) {

	opts := api.QueryOptions{
		WaitIndex: s.lastService,
		WaitTime:  5 * time.Minute,
	}
	if !reload {
		opts.WaitTime = time.Second
	}
	fmt.Println("Watching for new service with index", s.lastService, "or better")
	services, meta, err := catalog.Services(&opts)
	if err != nil {
		return
	}

	if meta.LastIndex > s.lastService {
		s.lastService = meta.LastIndex
	}

	myservices := make(map[string][]*api.CatalogService)
	for servicename := range services {
		// Get all instances for this service
		instances, _, _ := catalog.Service(servicename, "", nil)
		myservices[servicename] = instances
	}

	s.services = myservices
	s.buildConfig()
	if reload {
		reloadCaddy()
	}
}
