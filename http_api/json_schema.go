package http_api

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"github.com/complone/blast/common"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var jsonSchemaObj map[string]schemaInfo
var jsonSchemaRwlock sync.RWMutex

type schemaInfo struct {
	Schema *gojsonschema.Schema
}

func GetSchema(schemaFile string) *gojsonschema.Schema {
	jsonSchemaRwlock.RLock()
	defer jsonSchemaRwlock.RUnlock()
	common.Truef(jsonSchemaObj != nil, "jsonSchemaObj may not init correctly")
	if sh, ok := jsonSchemaObj[schemaFile]; ok {
		return sh.Schema
	}
	return nil
}

func ValidateJson(schema *gojsonschema.Schema, jsonBody []byte) error {
	documentLoader := gojsonschema.NewBytesLoader(jsonBody)
	rst, err := schema.Validate(documentLoader)
	if err != nil {
		return err
	}

	if rst.Valid() {
		return nil
	}

	var errtxt string
	for _, desc := range rst.Errors() {
		errtxt += fmt.Sprintf("%s(%s): %s\n", desc.Field(), desc.Value(), desc.Description())
	}

	return errors.New(errtxt)
}

func SyncJsonSchema(path string) error {
	if jsonSchemaObj == nil {
		jsonSchemaObj = make(map[string]schemaInfo)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Warningf("read dir %s failed. err = %s", path, err)
		return err
	}

	var newJsonSchemaObj = make(map[string]schemaInfo)
	for _, fi := range files {
		if !fi.IsDir() &&
			strings.HasPrefix(fi.Name(), "input.__") &&
			strings.HasSuffix(fi.Name(), ".schema.json") {
			prefix := strings.Split(fi.Name(), ".")[2]
			prefix = prefix[2 : len(prefix)-2]

			schemaStr, err := ioutil.ReadFile(path + "/" + fi.Name())
			if err != nil {
				log.Warningf("read file failed. err = %s", err)
				continue
			}

			loader := gojsonschema.NewStringLoader(string(schemaStr))
			schemaObj, err := gojsonschema.NewSchema(loader)
			if err != nil {
				log.Warningf("create schema failed. err = %s", err)
				continue
			}

			newJsonSchemaObj[prefix] = schemaInfo{schemaObj}
			log.Infof("load schema: %s -> %s", prefix, fi.Name())
		}
	}

	jsonSchemaRwlock.Lock()
	jsonSchemaObj = newJsonSchemaObj
	jsonSchemaRwlock.Unlock()

	return nil
}
