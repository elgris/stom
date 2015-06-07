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
)

// package settings
var (
	tagSetting          string      = "db"
	policySetting       Policy      = PolicyUseDefault
	defaultValueSetting interface{} = nil
)

type ToMappable interface {
	ToMap() (map[string]interface{}, error)
}

type ToMapper interface {
	ToMap(s interface{}) (map[string]interface{}, error)
}

type stom struct {
	defaultValue interface{}
	policy       Policy
	tag          string

	typ   reflect.Type
	cache map[string]int
}

func MustNewStom(s interface{}) *stom {
	typ, err := getStructType(s)
	if err != nil {
		panic(err.Error())
	}

	stom := &stom{
		typ:          typ,
		defaultValue: defaultValueSetting,
		policy:       policySetting,
	}
	stom.SetTag(tagSetting)

	return stom
}

func (this *stom) SetTag(tag string) {
	this.tag = tag
	this.cache = extractTagValues(this.typ, this.tag)
}
func (this *stom) SetDefault(defaultValue interface{}) { this.defaultValue = defaultValue }
func (this *stom) SetPolicy(policy Policy)             { this.policy = policy }

func (this *stom) ToMap(s interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		s = val.Elem().Interface()
	}

	typ := reflect.TypeOf(s)
	if typ != this.typ {
		return nil, errors.New(fmt.Sprintf("stom is set up to work with type %s, but %s given", this.typ, typ))
	}

	return toMap(s, this.cache, this.tag, this.defaultValue, this.policy)
}

func SetTag(t string)           { tagSetting = t }
func SetDefault(dv interface{}) { defaultValueSetting = dv }
func SetPolicy(p Policy)        { policySetting = p }

func ToMap(s interface{}) (map[string]interface{}, error) {
	if tomappable, ok := s.(ToMappable); ok {
		return tomappable.ToMap()
	}

	typ := reflect.TypeOf(s)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("expected struct, got %v", typ.Kind()))
	}

	tagmap := extractTagValues(typ, tagSetting)

	return toMap(s, tagmap, tagSetting, defaultValueSetting, policySetting)
}

func getStructType(s interface{}) (t reflect.Type, err error) {
	t = reflect.TypeOf(s)

	if t.Kind() == reflect.Invalid {
		err = errors.New(fmt.Sprintf("value is invalid:\n %v", s))
		return
	}

	if t.Kind() != reflect.Struct {
		err = errors.New(fmt.Sprintf("provided value is not a struct!\n%v", s))
	}

	return
}

func extractTagValues(typ reflect.Type, tag string) map[string]int {
	tagValues := map[string]int{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}

		if tagValue := field.Tag.Get(tagSetting); tagValue != "" && tagValue != "-" {
			tagValues[tagValue] = i
		}
	}

	return tagValues
}

func toMap(s interface{}, tagmap map[string]int, tag string, defaultValue interface{}, policy Policy) (map[string]interface{}, error) {
	val := reflect.ValueOf(s)

	result := map[string]interface{}{}

	for tag, index := range tagmap {
		vField := val.Field(index)

		v, err := convertValue(vField)
		if err != nil {
			return result, err
		}

		if v != nil {
			result[tag] = v
		} else if policy == PolicyUseDefault {
			result[tag] = defaultValue
		}

	}

	return result, nil
}

func convertValue(vField reflect.Value) (v interface{}, err error) {
	kind := vField.Kind()
	if kind == reflect.Ptr {
		if vField.Elem().IsValid() {
			v = vField.Elem().Interface()
		}
	} else {
		v = vField.Interface()
	}

	switch t := v.(type) {
	case driver.Valuer: // support for NullTypes like sql.NullString and so on
		if converted, convErr := t.Value(); convErr != nil || converted == nil {
			v = nil
		}
		return

	case ToMappable:
		v, err = t.ToMap()
		return
	}

	return v, nil

}
