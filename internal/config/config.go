package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	crplg "github.com/m1dugh/crawler/internal/plugin"

	yaml "gopkg.in/yaml.v2"
)

var ROOT_PATH string = func() string {

	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")

		if len(parts) > 1 && parts[0] == "HOME" {
			return parts[1] + "/.gocrawler"
		}
	}

	return ""
}()

var CONFIG_FILE = ROOT_PATH + "/config.yaml"

func initEmptyFile() {
	_, err := os.Open(CONFIG_FILE)
	if err != nil {
		_, err = os.Create(CONFIG_FILE)

		// root folder not created
		if err != nil {
			err = os.Mkdir(ROOT_PATH, 0777)
			if err != nil {
				log.Fatal(fmt.Sprintf("could not create root config repo for go crawler at %s", ROOT_PATH))
			}
			_, err = os.Create(CONFIG_FILE)
			if err != nil {
				log.Fatal(fmt.Sprintf("could not create config file at %s", CONFIG_FILE))
			}
		}
	}
}

func GetConfig() (Config, error) {

	initEmptyFile()
	source, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		return Config{}, errors.New("config::GetConfig -> could not read file")
	}

	var config Config
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return Config{}, errors.New("config::GetConfig -> could not unmarshal struct")
	}

	return config, nil
}

func SaveConfig(config Config) bool {

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return false
	}

	err = ioutil.WriteFile(CONFIG_FILE, bytes, 0777)
	return err == nil
}

func LoadPluginsFromConfig() map[string]*crplg.CrawlerPlugin {
	config, err := GetConfig()
	if err != nil {
		return make(map[string]*crplg.CrawlerPlugin, 0)
	}

	res := make(map[string]*crplg.CrawlerPlugin, len(config.Plugins))
	for _, pluginConfig := range config.Plugins {
		if pluginConfig.Active {
			// paths of plugins are relative to ROOT_PATH folder
			path := filepath.Join(ROOT_PATH, pluginConfig.Path)
			plg, err := crplg.GetCrawlerPlugin(path)
			if err == nil {
				res[pluginConfig.Name] = plg
			}
		}
	}

	return res
}