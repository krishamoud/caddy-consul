package caddyconsul

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

type domain struct {
	Config string
}

func (s *caddyfile) WatchKV(reload bool) {

	opts := api.QueryOptions{
		WaitIndex: s.lastKV,
		WaitTime:  5 * time.Minute,
	}
	if !reload {
		opts.WaitTime = time.Second
	}
	fmt.Println("Watching for new KV with index", s.lastKV, "or better")
	pairs, meta, err := kv.List("caddy/", &opts)
	if err != nil {
		fmt.Println(err)
		// this should probably be logged
		return
	}
	if meta.LastIndex > s.lastKV {
		s.lastKV = meta.LastIndex
	}
	// If there's nothing, at least put our KV value so the user isn't lost
	if len(pairs) == 0 {
		kv.Put(&api.KVPair{Key: "caddy/"}, nil)
	}

	// TODO actually make a new one, don't just keep using the old one
	domains := make(map[string]*domain)
	for _, k := range pairs {
		keybits := strings.SplitN(k.Key, "/", 3)
		if len(keybits) != 2 || keybits[1] == "" {
			continue
		}
		domains[keybits[1]] = &domain{
			Config: string(k.Value),
		}
	}
	s.domains = domains
	s.buildConfig()

	if reload {
		reloadCaddy()
	}
}
