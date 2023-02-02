package vmetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Message struct {
	messages []interface{}
}

func NewMessage() *Message {
	return &Message{
		messages: make([]interface{}, 0),
	}
}

func toMap(intf interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := json.Marshal(&intf)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(b, &m)
	return m, err
}

func (m *Message) jsonString(inf interface{}) (string, error) {
	kind := reflect.ValueOf(inf).Kind()
	if kind != reflect.Struct && kind != reflect.Ptr {
		return "", errors.New("message should be type of struct: " + reflect.ValueOf(inf).Kind().String())
	}
	infMap, err := toMap(inf)
	if err != nil {
		return "", errors.New("message should be type of struct: " + reflect.ValueOf(inf).Kind().String())
	}
	infMap["__time__"] = time.Now().UnixNano()
	marshal, err := json.Marshal(infMap)
	if err != nil {
		return "", err
	}
	return string(marshal), err
}

func (m *Message) Record(inf interface{}) {
	m.messages = append(m.messages, inf)
}

func (m *Message) Consume() []string {
	var messages []string
	for _, msg := range m.messages {
		jsonString, err := m.jsonString(msg)
		if err != nil {
			fmt.Println(err)
		} else {
			messages = append(messages, jsonString)
		}
	}
	return messages
}

func (m *Message) Clear() {
	m.messages = nil
	m.messages = make([]interface{}, 0)
}
