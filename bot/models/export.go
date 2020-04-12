package models

import (
	"reflect"
	"strings"

	"github.com/fatih/structtag"
)

func Export(input interface{}) map[string]interface{} {
	var (
		kind  = reflect.TypeOf(input)
		value = reflect.ValueOf(input)
	)

	res := map[string]interface{}{}
	for i := 0; i < kind.NumField(); i++ {
		tags, err := structtag.Parse(string(kind.Field(i).Tag))
		if err != nil {
			continue
		}
		tag, err := tags.Get("db")
		if err != nil {
			continue
		}
		res[tag.Name] = value.Field(i).Interface()
	}
	return res
}

func Returning(input interface{}) string {
	fields := []string{}
	kind := reflect.TypeOf(input)
	for i := 0; i < kind.NumField(); i++ {
		field := kind.Field(i)
		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			continue
		}
		name, err := tags.Get("db")
		if err != nil {
			continue
		}
		fields = append(fields, name.Name)
	}
	if len(fields) == 0 {
		return ""
	}
	return "RETURNING " + strings.Join(fields, ", ")
}
