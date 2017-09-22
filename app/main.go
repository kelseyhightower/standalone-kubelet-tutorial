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

var version = "0.1.0"

type Config struct {
	Hostname string `json:"hostname"`
	Key      int64  `json:"key"`
}

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

		fmt.Fprintf(w, "version: %s\nhostname: %s\nkey: %s\n", version, config.Hostname, config.Key)
	})

	log.Println("Starting the HTTP service...")
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
