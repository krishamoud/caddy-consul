package caddyconsul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"strings"
	"time"
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
	ret := ""
	for domain, services := range s.services {
		// policy := ""
		// websockets := ""
		// // var port string
		// if ok := contains(services[0].ServiceTags, "websockets"); ok {
		// 	websockets = "websocket"
		// }
		// if ok := contains(services[0].ServiceTags, "stickysessions"); ok {
		// 	policy = "ip_hash"
		// } else {
		// 	policy = "random"
		// }

		// if len(services[0].ServiceTags) > 4 {
		// 	explicitPort := services[0].ServiceTags[4]
		// 	if _, err := strconv.Atoi(explicitPort); err == nil {
		// 		port = explicitPort
		// 	} else {
		// 		port = strconv.Itoa(services[0].ServicePort)
		// 	}
		// } else {
		// 	port = strconv.Itoa(services[0].ServicePort)
		// }
		alias := strings.Contains(domain, "rompliapp.com")
		if alias {
			ret += "proxy " + domain + ":80 " + services[0].Address + ":" + strconv.Itoa(services[0].ServicePort) + " {\n"
			// ret += domain + " {\n"
			// ret += "  tls /ssl/ssl-bundle.crt /ssl/star_rompliapp_com.key\n"
			ret += "  tls off\n"
			// ret += "  redir 301 {\n"
			// ret += "    if {>X-Forwarded-Proto} is http\n"
			// ret += "    /  https://{host}{uri}\n"
			ret += "  }\n"
		} else {
			ret += domain + " {\n"
			ret += "  tls {\n"
			ret += "    max_certs 10\n"
			ret += "  }\n"
		}
		// ret += "  gzip\n"
		// ret += "  proxy / { \n"
		// ret += "    policy " + policy + "\n"
		// if len(websockets) > 0 {
		// 	ret += "    websocket \n"
		// }
		// ret += "    transparent \n"
		// // ret += "    " + websockets + " \n"
		// for _, service := range services {
		// 	// fmt.Println(service)
		// 	ret += "    upstream " + service.Address + ":" + strconv.Itoa(service.ServicePort) + "\n"
		// }
		// ret += "  } \n"
		// ret += "}\n" // close domain
		// ret += "http://" + domain + " {\n"
		// ret += "  redir 301 {\n"
		// ret += "    /  https://{host}{uri}\n"
		// ret += "  }\n"
		// ret += "}\n" // close domain
	}
	// ret += "tunguska-gauge-demo.rompli.com {\n"
	// ret += "   redir 301 {\n"
	// ret += "    /  https://tunguska-gauge-demo.rompliapp.com{uri}\n"
	// ret += "  }\n"
	// ret += "}\n"
	// ret += "league.rompli.com {\n"
	// ret += "   redir 301 {\n"
	// ret += "    /  https://league.rompliapp.com{uri}\n"
	// ret += "  }\n"
	// ret += "}\n"
	// ret += "aztbc.rompli.com {\n"
	// ret += "   redir 301 {\n"
	// ret += "    /  https://aztbc.rompliapp.com{uri}\n"
	// ret += "  }\n"
	// ret += "}\n"
	s.contents = ret
}
