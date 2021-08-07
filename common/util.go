package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func TrimArray(arr []string) {
	for i, v := range arr {
		arr[i] = strings.Trim(v, " ")
	}
}

func Interface2int(v interface{}) (int64, error) {
	switch value := v.(type) { // original type
	case string:
		return strconv.ParseInt(value, 10, 64)
	case json.Number:
		return v.(json.Number).Int64()
	case int32:
		return int64(value), nil
	case uint32:
		return int64(value), nil
	case int:
		return int64(value), nil
	case uint:
		return int64(value), nil
	case float32:
		return int64(value), nil
	case float64:
		return int64(value), nil
	case int64:
		return value, nil
	case uint64:
		return int64(value), nil
	default:
		return 0, errors.New(fmt.Sprintf("unknown type: %T", value))
	}
}

func Interface2float64(v interface{}) (float64, error) {
	switch value := v.(type) { // original type
	case string:
		return strconv.ParseFloat(value, 64)
	case json.Number:
		return v.(json.Number).Float64()
	case int32:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case int:
		return float64(value), nil
	case uint:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case int64:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	default:
		return 0.0, errors.New(fmt.Sprintf("unknown type: %T", value))
	}
}

type Condition struct {
	data []string
}

func NewCondition() *Condition {
	return &Condition{
		data: make([]string, 0),
	}
}

func (c *Condition) AddCond(cond string) {
	c.data = append(c.data, cond)
}

func (c *Condition) Dump(prefix string, conj string) string {
	if len(c.data) == 0 {
		return ""
	}

	return " " + prefix + " " + strings.Join(c.data, " "+conj+" ")
}

func IntArrayToStrArray(arr []int64) []string {
	idStrs := make([]string, 0)
	for _, id := range arr {
		idStrs = append(idStrs, strconv.FormatInt(id, 10))
	}

	return idStrs
}

func GetStringFromMap(m *map[string]interface{}, key string, def string) string {
	if v, ok := (*m)[key]; ok {
		return v.(string)
	}

	return def
}

func GetBoolFromMap(m *map[string]interface{}, key string, def bool) bool {
	if v, ok := (*m)[key]; ok {
		return v.(bool)
	}

	return def
}

func GetIntFromMap(m *map[string]interface{}, key string, def int) int {
	if v, ok := (*m)[key]; ok {
		nv, err := Interface2int(v)
		Truef(err == nil, "err = %s", err)
		return int(nv)
	}

	return def
}

func QuoteJoin(strs []string, sep string) string {
	for i, item := range strs {
		strs[i] = strconv.Quote(item)
	}

	return strings.Join(strs, sep)
}

