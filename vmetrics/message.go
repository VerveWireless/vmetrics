package vmetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Message struct {
	Namespace string
	Subsystem string
	Name       string
	messages []interface{}
}

func NewMessage(namespace, subsystem, name string) *Message {
	return &Message{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:       name,
		messages: make([]interface{}, 0),
	}
}

func toMap(intf interface{}) map[string]interface{} {
	b, _ := json.Marshal(&intf)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

func (m *Message) jsonString(inf interface{}) (string, error) {
	if reflect.ValueOf(inf).Kind() != reflect.Struct {
		return "", errors.New("message should be type of struct")
	}
	infMap := toMap(inf)
	infMap["__time__"] = time.Now().UnixNano()
	infMap["__namespace__"] = m.Namespace
	infMap["__subsystem__"] = m.Subsystem
	infMap["__name__"] = m.Name
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

func (m *Message) GetName() string {
	return m.Name
}