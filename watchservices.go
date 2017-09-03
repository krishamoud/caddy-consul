package caddyconsul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type service struct {
	Name           string
	Instances      []*api.CatalogService
	Websockets     bool
	StickySessions bool
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
	// TODO
	services, meta, err := catalog.Services(&opts)
	if err != nil {
		// TODO should probably handle this better
		return
	}

	if meta.LastIndex > s.lastService {
		s.lastService = meta.LastIndex
	}

	myservices := make(map[string][]*api.CatalogService)
	for servicename := range services {
		// Get all instances for this service
		instances, _, _ := catalog.Service(servicename, "", nil)
		// TODO should probably check error
		for _, instance := range instances {
			if ok := contains(instance.ServiceTags, "udp"); ok {
				fmt.Println(instance)
				myservices[instance.ServiceTags[0]+".rompliapp.com"] = instances
			}
			// if len(instance.ServiceTags) >= 4 {
			// 	myservices[instance.ServiceTags[0]+".rompliapp.com"] = instances
			// 	if instance.ServiceTags[0] != instance.ServiceTags[1] {
			// 		myservices[instance.ServiceTags[1]] = instances
			// 	}
			// }
		}
	}
	fmt.Println("Services:", myservices)
	s.services = myservices

	s.buildConfig()

	if reload {
		reloadCaddy()
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
