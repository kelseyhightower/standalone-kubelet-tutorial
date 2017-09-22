package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	config  string
	onetime bool
)

type Config struct {
	Hostname string `json:"hostname"`
	Key      int64  `json:"key"`
}

func main() {
	flag.StringVar(&config, "config", "", "Path to config file.")
	flag.BoolVar(&onetime, "onetime", false, "Generate config and exit")
	flag.Parse()

	if config == "" {
		log.Fatal("The -config flag is required. Exiting...")
	}

	for {
		log.Println("Generating configuration file...")
		c, err := getConfig()

		if err != nil {
			time.Sleep(30 * time.Second)
			log.Println(err)
			continue
		}

		data, err := json.MarshalIndent(&c, "", " ")
		if err != nil {
			time.Sleep(30 * time.Second)
			log.Println(err)
			continue
		}

		if err := ioutil.WriteFile(config, data, 0644); err != nil {
			time.Sleep(30 * time.Second)
			log.Println(err)
			continue
		}

		log.Printf("Wrote config %s", config)

		if onetime {
			log.Println("Onetime mode set. Exiting...")
			os.Exit(0)
		}

		time.Sleep(30 * time.Second)
	}
}

func getConfig() (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	c := Config{
		Hostname: hostname,
		Key:      time.Now().Unix(),
	}

	return &c, nil
}
