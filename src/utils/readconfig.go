package utils

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	//"sync"
)


type mainConfig struct {
	Scheduler SchInfo `json:"scheduler"`
	Workers []WorkerNodeInfo `json:"workers"`
	Servers  []ServerNodeInfo `json:"servers"`
}

type SchInfo struct {
	Address string `json:"address"`
	Consistency string `json:"consistency"`
}

type WorkerNodeInfo struct {
	ID string `json:"index"`
	Address string `json:"address"`
	Prepf string `json:"prepf"`
}

type ServerNodeInfo struct {
	ID string `json:"index"`
	Address string `json:"address"`
	Prepf string `json:"prepf"`
}


var globalConfigPath string = "../../configs/godml.json"

func LoadGlobalConfig() *mainConfig {
	return loadConfig(globalConfigPath)
}


func loadConfig(path string) *mainConfig{
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config conf failed: ", err)
	}
	mainConfig := &mainConfig{}
	err = json.Unmarshal(buf, mainConfig)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	//fmt.Printf("%+v \n", mainConfig)
	return mainConfig
}


// func Init(path string) *MainConfig {
// 	mainConfig := LoadConfig(path)
// 	conf = mainConfig
// 	return conf
// }



