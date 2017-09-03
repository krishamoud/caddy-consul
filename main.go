package caddyconsul

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/mholt/caddy"
)

var consulClient *api.Client
var kv *api.KV
var catalog *api.Catalog
var caddyfilePath string

var started = time.Now()

func init() {
	caddy.RegisterCaddyfileLoader("myloader", caddy.LoaderFunc(myLoader))
}

func reloadCaddy() {
	if time.Since(started) < 2*time.Second {
		fmt.Println("Not reloading since caddy uptime is too short")
		return
	}
	self, _ := os.FindProcess(os.Getpid())
	self.Signal(syscall.SIGUSR1)
}

func myLoader(serverType string) (caddy.Input, error) {

	// return early, if we can
	if consulGenerator != nil {
		return consulGenerator, nil
	}

	// Assume localhost, if it's not set in the environment
	consulAddress := os.Getenv("CONSUL")
	caddyfilePath = os.Getenv("CADDY_FILE_PATH")
	if consulAddress == "" {
		consulAddress = "127.0.0.1:8500"
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress

	var err error

	if consulClient == nil {
		// setup our consulClient connection
		consulClient, err = api.NewClient(consulConfig)
		if err != nil {
			return nil, err
		}
	}

	// setup our KV connection
	kv = consulClient.KV()
	// setup our catalog connection
	catalog = consulClient.Catalog()

	// Actually create the right instance as a generator that caddy needs
	consulGenerator = new(caddyfile)
	// let the KV and Service portions generate once so we have content for the caddy file when caddy asks the first time
	consulGenerator.WatchServices(false)

	if len(consulGenerator.Body()) == 0 {
		return nil, nil
	}
	// Start our loop that keeps checking on consul
	consulGenerator.StartWatching()
	return consulGenerator, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
