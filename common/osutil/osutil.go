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
	log.Println("Creating working directory for " + name)

	mrPath := os.Getenv("MR_PATH")
	if mrPath == "" {
		panic("Set MR_PATH to specify file system data location")
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
			log.Println("Couldn't load " + configName + ", error: " + err.Error())
			log.Println("Config loaded from CLI arguments")
		} else {
			log.Println("Config loaded from " + configName)
		}
	} else {
		log.Println("Config overriden from CLI arguments")
	}
	if err := SaveConfig(config); err != nil {
		log.Println("Couldn't save config to " + configName + ", error: " + err.Error())
	} else {
		log.Println("Saved config to " + configName)
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

func SaveConfig(config interface{}) error {
	file, err := os.Create(configName)
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
