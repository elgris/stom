package stom

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
)

type Policy uint8

const (
	// PolicyUseDefault enforces SToM to use defined default value instead
	// of 'nil' value in resulting map[string]interface{}.
	// Default value can be nil also.
	PolicyUseDefault Policy = iota

	// PolicyExclude tells SToM to ignore 'nil' values and to not include them in
	// resulting map
	PolicyExclude
)

// Package settings
// They are used as defaults for initialization if new SToMs
var (
	tagSetting          string      = "db"
	policySetting       Policy      = PolicyUseDefault
	defaultValueSetting interface{} = nil
)

// ToMappable defines an entity that knows how to convert itself to map[string]interface{}.
// If an entity implements this interface, SToM won't do any magic,
// it will just call ToMap() method to makes thigs simpler.
// Such approach allows custom conversion
type ToMappable interface {
	ToMap() (map[string]interface{}, error)
}

// ToMapper defines a service that is able to convert something to map[string]interface{}
type ToMapper interface {
	ToMap(s interface{}) (map[string]interface{}, error)
}

// stom is a small handy tool that is instantiated for certain type and caches
// all knowledge about this type to increase conversion speed
type stom struct {
	defaultValue interface{}
	policy       Policy
	tag          string

	typ   reflect.Type
	cache map[string]int
}

// MustNewStom creates new instance of a SToM converter for type of given structure.
// Panics if no structure provided
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

// SetTag sets SToM to scan for given tag in structure
func (this *stom) SetTag(tag string) {
	this.tag = tag
	this.cache = extractTagValues(this.typ, this.tag)
}

// SetDefault makes SToM to put given default value in 'nil' values of structure's fields
func (this *stom) SetDefault(defaultValue interface{}) { this.defaultValue = defaultValue }

// SetPolicy sets policy for 'nil' values
func (this *stom) SetPolicy(policy Policy) { this.policy = policy }

// ToMap converts a structure to map[string]interface{}.
// SToM converts only structures it was initialized for
func (this *stom) ToMap(s interface{}) (map[string]interface{}, error) {
	typ, err := getStructType(s)
	if err != nil {
		return nil, err
	}

	if typ != this.typ {
		return nil, errors.New(fmt.Sprintf("stom is set up to work with type %s, but %s given", this.typ, typ))
	}

	return toMap(s, this.cache, this.tag, this.defaultValue, this.policy)
}

// SetTag sets package setting for tag to look for in incoming structures
func SetTag(t string) { tagSetting = t }

// SetDefault sets package default value to set instead of 'nil' in resulting maps
func SetDefault(dv interface{}) { defaultValueSetting = dv }

// SetPolicy sets package setting for policy. Policy defines what to do with
// 'nil' values in resulting maps.
// There are 2 policies:
// - PolicyUseDefault - with this policy default value will be used instead of 'nil'
// - PolicyExclude    - 'nil' values will be discarded
func SetPolicy(p Policy) { policySetting = p }

// ToMap converts given structure into map[string]interface{}
func ToMap(s interface{}) (map[string]interface{}, error) {
	if tomappable, ok := s.(ToMappable); ok {
		return tomappable.ToMap()
	}

	typ, err := getStructType(s)
	if err != nil {
		return nil, err
	}

	tagmap := extractTagValues(typ, tagSetting)

	return toMap(s, tagmap, tagSetting, defaultValueSetting, policySetting)
}

func getStructType(s interface{}) (t reflect.Type, err error) {
	t = reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Invalid {
		err = errors.New(fmt.Sprintf("value is invalid:\n %v", s))
		return
	}

	if t.Kind() != reflect.Struct {
		err = errors.New(fmt.Sprintf("provided value is not a struct but %v!", t.Kind()))
	}

	return
}

// extractTagValues scans given type and tries to find all fields with given tag
// Indices of all found fields are stored as values in resulting map
// Keys of resulting map are actual values of tags
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

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := map[string]interface{}{}

	for tag, index := range tagmap {
		vField := val.Field(index)

		v, err := filterValue(vField)
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

// filterValue filters given value of some structure's field.
// Simple values are left as is. There is some special logic about particular types
// like ToMappable
func filterValue(vField reflect.Value) (v interface{}, err error) {
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
