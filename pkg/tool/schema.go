package tool

import (
	"fmt"
	"reflect"
	"strings"
)

func StructToJSONSchema(data any) (map[string]any, error) {
	t := reflect.TypeOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	schema := map[string]any{
		"type":       "object",
		"required":   []any{},
		"properties": map[string]any{},
	}

	properties := schema["properties"].(map[string]any)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonTag := field.Tag.Get("json")
		jsonTagParts := strings.Split(jsonTag, ",")
		jsonTagName := jsonTagParts[0]

		if jsonTagName == "-" || jsonTagName == "" {
			continue
		}

		if len(jsonTagParts) > 1 && jsonTagParts[1] == "required" {
			schema["required"] = append(schema["required"].([]any), jsonTagName)
		}

		fieldType := field.Type.Kind()

		switch fieldType {
		case reflect.String:
			properties[jsonTagName] = map[string]any{
				"type": "string",
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			properties[jsonTagName] = map[string]any{
				"type": "integer",
			}
		case reflect.Float32, reflect.Float64:
			properties[jsonTagName] = map[string]any{
				"type": "number",
			}
		case reflect.Bool:
			properties[jsonTagName] = map[string]any{
				"type": "boolean",
			}
		case reflect.Struct:
			nestedSchema, err := StructToJSONSchema(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
			properties[jsonTagName] = nestedSchema
		case reflect.Slice:
			elemType := field.Type.Elem().Kind()
			var itemsType string
			switch elemType {
			case reflect.String:
				itemsType = "string"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				itemsType = "integer"
			case reflect.Float32, reflect.Float64:
				itemsType = "number"
			case reflect.Bool:
				itemsType = "boolean"
			case reflect.Struct:
				nestedSchema, err := StructToJSONSchema(reflect.New(field.Type.Elem()).Elem().Interface())
				if err != nil {
					return nil, err
				}
				properties[jsonTagName] = map[string]any{
					"type":  "array",
					"items": nestedSchema,
				}
				continue
			default:
				itemsType = "object"
			}
			properties[jsonTagName] = map[string]any{
				"type":  "array",
				"items": map[string]any{"type": itemsType},
			}
		default:
			properties[jsonTagName] = map[string]any{
				"type": "object",
			}
		}
	}

	return schema, nil
}
