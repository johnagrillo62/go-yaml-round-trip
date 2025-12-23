package nested

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

// ValidateYAMLMatchesStruct validates a YAML file against any struct recursively
func ValidateYAMLMatchesStruct(yamlFile string, structType interface{}) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read YAML: %w", err)
	}

	var yamlMap map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlMap); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return validateMap(yamlMap, reflect.TypeOf(structType))
}

// validateMap recursively compares a map to a struct type
func validateMap(yamlMap map[string]interface{}, t reflect.Type) error {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", t.Kind())
	}

	// Check all struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}

		value, exists := yamlMap[tag]
		if !exists {
			return fmt.Errorf("missing field '%s' in YAML", tag)
		}

		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		switch ft.Kind() {
		case reflect.Struct:
			subMap, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("field '%s' should be a map in YAML", tag)
			}
			if err := validateMap(subMap, ft); err != nil {
				return err
			}
		case reflect.Slice:
			slice, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("field '%s' should be a slice in YAML", tag)
			}
			elemType := ft.Elem()
			for _, item := range slice {
				switch elemType.Kind() {
				case reflect.Struct:
					itemMap, ok := item.(map[string]interface{})
					if !ok {
						return fmt.Errorf("slice element in '%s' should be a map", tag)
					}
					if err := validateMap(itemMap, elemType); err != nil {
						return err
					}
				default:
					// primitive slice; no further validation
				}
			}
		default:
			// primitive type; nothing to recurse
		}
	}

	// Check for extra keys in YAML that don't exist in struct
	for key := range yamlMap {
		found := false
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("json")
			if tag == "" {
				tag = field.Name
			}
			if key == tag {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("YAML key '%s' does not exist in struct %s", key, t.Name())
		}
	}

	return nil
}
