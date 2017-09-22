package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addr       string
	configFile string
)

var version = "0.2.0"

type Config struct {
	Hostname string `json:"hostname"`
	Key      int64  `json:"key"`
}

var indexPage = `<!doctype html>
<html lang="en">
<head>
  <title>Standalone Kubelet</title>
</head>

<body>
  <h1>App</h1>
  <h3>Version</h3>
  <ul>
    <li>%s</li>
  </ul>

  <h3>Config</h3>
  <ul>
    <li>Hostname: %s</li>
    <li>Key: %d</li>
  </ul>
</body>
</html>
`

func main() {
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "HTTP listen address.")
	flag.StringVar(&configFile, "config", "/etc/app/config.json", "Path to config file.")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		config, err := getConfig(configFile)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		fmt.Fprintf(w, indexPage, version, config.Hostname, config.Key)
	})

	http.ListenAndServe(addr, nil)
}

func getConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
