package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	OpenAPI struct {
		HttpAddr   string `json:"http_addr"`
		APISwitch  string `json:"api_switch"` // query,upserts
		Clutser    string `json:"cluster"`
		SchemaPath string `json:"schema_path"`
	} `json:"openapi"`
}

type IConfig interface {
	GetConfig() Config
}

var GlobalConfig Config

func InitConfigFromJson(jsonFile string, cfg IConfig) error {
	f, err := os.Open(jsonFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, cfg)
	if err != nil {
		return err
	}

	GlobalConfig = cfg.GetConfig()

	if GlobalConfig.OpenAPI.HttpAddr == "" {
		return errors.New("openapi.http_addr is required")
	}

	if GlobalConfig.OpenAPI.APISwitch == "" {
		return errors.New("openapi.api_switch is required")
	}

	return nil
}
