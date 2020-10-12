package config

import (
    "errors"
    "fmt"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
)

type Config struct {
    Mysql   mysql   `yaml:"mysql"`
    Redis   redis   `yaml:"redis"`
    Session session `yaml:"session"`
    Auth    auth    `yaml:"auth"`
    Aes     aes     `yaml:"aes"`
}

type mysql struct {
    Dsn   string `yaml:"dsn"`
    Debug int32  `yaml:"debug"`
}

type redis struct {
    Host     string `yaml:"host"`
    Port     string `yaml:"port"`
    Password string `yaml:"password"`
    Db       int32  `yaml:"db"`
}

type session struct {
    Host     string `yaml:"host"`
    Port     string `yaml:"port"`
    Password string `yaml:"password"`
    Db       int32  `yaml:"db"`
}

type auth struct {
    Skip string `yaml:"skip"`
}

type aes struct {
    Key string `yaml:"key"`
}

var (
    Cfg Config
)

func LoadCfg(env string) (err error) {
    if env == "" {
        err = errors.New("load config with env error")
        logger.Logger.Errorf("env error: %v", err)
        return
    }

    dir, _ := os.Getwd()
    var fs []byte
    fs, err = ioutil.ReadFile(fmt.Sprintf("%s/config/%s.yml", dir, env))
    if err != nil {
        logger.Logger.Errorf("read config file error: %v", err)
        return
    }

    conf := Config{}
    err = yaml.Unmarshal(fs, &conf)
    if err != nil {
        logger.Logger.Errorf("load cfg with unmarshal error: %v", err)
        return
    }

    Cfg = conf
    return
}
