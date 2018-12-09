package osutil

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const configName = "config.json"

func Init(name string, override bool, config interface{}) {
	log.Println("creating working directory for " + name)

	mrPath := os.Getenv("MR_PATH")
	if mrPath == "" {
		panic("set MR_PATH to specify file system data location")
	}

	workingDir := path.Join(mrPath, name)
	if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.Chdir(workingDir); err != nil {
		panic(err)
	}

	if !override {
		if err := getConfig(config); err != nil {
			log.Println("couldn't load " + configName + ", error: " + err.Error())
			log.Println("config loaded from CLI arguments")
		} else {
			log.Println("config loaded from " + configName)
		}
	} else {
		log.Println("config overriden from CLI arguments")
	}
	if err := saveConfig(config); err != nil {
		log.Println("couldn't save config to " + configName + ", error: " + err.Error())
	} else {
		log.Println("saved config to " + configName)
	}
}

func getConfig(config interface{}) error {
	_, err := os.Stat(configName)
	if os.IsNotExist(err) {
		return err
	}
	file, err := os.Open(configName)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, config)
}

func saveConfig(config interface{}) error {
	file, err := os.OpenFile(configName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err
}
