package vmetrics

import (
	"encoding/json"
	"errors"
	"reflect"
)

func toMap(intf interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	kind := reflect.ValueOf(intf).Kind()
	if kind != reflect.Struct && kind != reflect.Ptr {
		return m, errors.New("message should be type of struct: " + reflect.ValueOf(intf).Kind().String())
	}
	b, err := json.Marshal(&intf)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(b, &m)
	return m, err
}

func jsonString(inf interface{}) (string, error) {
	infMap, err := toMap(inf)
	if err != nil {
		return "", errors.New("message should be type of struct: " + reflect.ValueOf(inf).Kind().String())
	}
	//infMap["__time__"] = time.Now().UnixNano()
	marshal, err := json.Marshal(infMap)
	if err != nil {
		return "", err
	}
	return string(marshal), err
}
