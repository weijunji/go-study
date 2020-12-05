package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/go-homedir"
	"github.com/qiniu/qmgo"
	"gopkg.in/yaml.v2"
)

func getConfigFile() []byte {
	path, err := homedir.Expand("~/lottery_conf.yaml")
	if err != nil {
		panic("get homedir failed")
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic("read file failed")
	}
	return file
}

// GetConfig get config in lottery_conf.yaml
func GetConfig(namespace string) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(getConfigFile(), &m)
	if err != nil {
		panic(err)
	}
	return m[namespace].(map[interface{}]interface{})
}

func getMysqlSource() string {
	config := GetConfig("mysql")
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["user"],
		config["password"],
		config["host"],
		config["port"],
		config["database"],
	)
}

func getRedisOption() *redis.Options {
	config := GetConfig("redis")
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config["host"], config["port"]),
		Password: config["password"].(string),
		DB:       config["database"].(int),
	}
}

func getMongoConfig() *qmgo.Config {
	config := GetConfig("mongodb")
	return &qmgo.Config{
		Uri:      fmt.Sprintf("mongodb://%s:%d", config["host"], config["port"]),
		Database: config["database"].(string),
		Auth: &qmgo.Credential{
			Username: config["user"].(string),
			Password: config["password"].(string),
		},
	}
}
