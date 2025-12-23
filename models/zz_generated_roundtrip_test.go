package main

import (
	"os"
	"reflect"
	"testing"
	"gopkg.in/yaml.v3"
)

func validateYAMLMatchesStruct(t *testing.T, yamlFile string, structType interface{}) {
	data, err := os.ReadFile(yamlFile)
	if err != nil { t.Fatal(err) }

	var yamlMap map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlMap); err != nil { t.Fatal(err) }

	val := reflect.TypeOf(structType)
	if val.Kind() == reflect.Ptr { val = val.Elem() }

	tags := map[string]bool{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" { tag = field.Name }
		tags[tag] = true
	}

	for key := range yamlMap {
		if !tags[key] {
			t.Errorf("YAML key '%s' does not exist in struct %s", key, val.Name())
		}
	}

	for tag := range tags {
		if _, ok := yamlMap[tag]; !ok {
			t.Errorf("Struct field '%s' is missing in YAML", tag)
		}
	}
}


func TestYAML_User_Valid(t *testing.T) {
	validateYAMLMatchesStruct(t, "valid_user.yaml", User{})
}

func TestYAML_User_Invalid(t *testing.T) {
	validateYAMLMatchesStruct(t, "invalid_user.yaml", User{})
}

func TestYAML_Car_Valid(t *testing.T) {
	validateYAMLMatchesStruct(t, "valid_car.yaml", Car{})
}

func TestYAML_Car_Invalid(t *testing.T) {
	validateYAMLMatchesStruct(t, "invalid_car.yaml", Car{})
}

