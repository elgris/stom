package stom

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
)

type Policy uint8

const (
	PolicyUseDefault Policy = iota
	PolicyExclude

	defaultTag = "db"
)

type ToMappable interface {
	ToMap() (map[string]interface{}, error)
}

var defaultStom = stom{defaultTag, PolicyUseDefault, nil}

type stom struct {
	tag          string
	policy       Policy
	defaultValue interface{}
}

func (this *stom) SetTag(tag string) {
	this.tag = tag
}

func (this *stom) SetDefault(defaultValue interface{}) {
	this.defaultValue = defaultValue
}

func (this *stom) SetPolicy(policy Policy) {
	this.policy = policy
}

func (this *stom) ToMap(s interface{}) (map[string]interface{}, error) {
	typ := reflect.TypeOf(s)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("expected struct, got %v", typ.Kind()))
	}

	result := map[string]interface{}{}

	val := reflect.ValueOf(s)
	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}

		// TODO: try indices instead of values
		// TODO: parallelize it
		// TODO: cache should work as so:
		// 1. Create a map 'field name' => index
		// 2. Get a value by index
		// 3. Convert value to appropriate value
		// 4. DONE

		if t := field.Tag.Get(this.tag); t != "" && t != "-" {
			vField := val.Field(i)
			var v interface{}
			if vField.Kind() == reflect.Ptr {
				if vField.Elem().IsValid() {
					v = vField.Elem().Interface()
				}
			} else {
				v = vField.Interface()
			}

			v = convertValue(v)

			if v != nil {
				result[t] = v
			} else if this.policy == PolicyUseDefault {
				result[t] = this.defaultValue
			}
		}

	}

	return result, nil

	// TODO:
	// 1. Check if it's a struct. If not - return error
	// 2. Scan through struct's fields and get all the tags
	// 3. If tag is in this.tags - use tag name as a key and field value as value
	// and put them into map
}

func SetDefault(defaultValue interface{})                 { defaultStom.SetDefault(defaultValue) }
func SetTag(tag string)                                   { defaultStom.SetTag(tag) }
func SetPolicy(policy Policy)                             { defaultStom.SetPolicy(policy) }
func ToMap(s interface{}) (map[string]interface{}, error) { return defaultStom.ToMap(s) }

func convertValue(input interface{}) (output interface{}) {
	// TODO: check if input is a structure
	output = input
	switch t := input.(type) {
	case driver.Valuer:
		if converted, err := t.Value(); converted == nil || err != nil {
			output = nil
		}
	}

	return

}
