// Package stom is about converting structs to map[string]interface{} with
// minimum processing and overhead
package stom

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// Policy is a type to define policy of dealing with 'nil' values
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
	tagSetting          = "db"
	policySetting       = PolicyUseDefault
	defaultValueSetting interface{}
)

// Zeroable is an interface that allows to filter values that can explicitly
// state that they are 'zeroes'. For example, this interface allows to filter
// zero time.Time,
type Zeroable interface {
	IsZero() bool
}

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

// ToMapperFunc defines a function that implements ToMapper
type ToMapperFunc func(s interface{}) (map[string]interface{}, error)

// ToMap implements ToMapper
func (f ToMapperFunc) ToMap(s interface{}) (map[string]interface{}, error) {
	return f(s)
}

type tags struct {
	Simple map[int]string
	Nested map[int]tags
}

func newTags() tags {
	return tags{
		Simple: make(map[int]string),
		Nested: make(map[int]tags),
	}
}

// stom is a small handy tool that is instantiated for certain type and caches
// all knowledge about this type to increase conversion speed
type stom struct {
	defaultValue interface{}
	policy       Policy
	tag          string

	typ   reflect.Type
	cache tags
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
func (s *stom) SetTag(tag string) *stom {
	s.tag = tag
	s.cache = extractTagValues(s.typ, s.tag)

	return s
}

// SetDefault makes SToM to put given default value in 'nil' values of structure's fields
func (s *stom) SetDefault(defaultValue interface{}) *stom {
	s.defaultValue = defaultValue

	return s
}

// SetPolicy sets policy for 'nil' values
func (s *stom) SetPolicy(policy Policy) *stom {
	s.policy = policy

	return s
}

// ToMap converts a structure to map[string]interface{}.
// SToM converts only structures it was initialized for
func (s *stom) ToMap(obj interface{}) (map[string]interface{}, error) {
	typ, err := getStructType(obj)
	if err != nil {
		return nil, err
	}

	if typ != s.typ {
		return nil, fmt.Errorf("stom is set up to work with type %s, but %s given", s.typ, typ)
	}

	return toMap(obj, s.cache, s.defaultValue, s.policy)
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

// ConvertToMap converts given structure into map[string]interface{}
func ConvertToMap(s interface{}) (map[string]interface{}, error) {
	if tomappable, ok := s.(ToMappable); ok {
		return tomappable.ToMap()
	}

	typ, err := getStructType(s)
	if err != nil {
		return nil, err
	}

	tagmap := extractTagValues(typ, tagSetting)

	return toMap(s, tagmap, defaultValueSetting, policySetting)
}

func getStructType(s interface{}) (t reflect.Type, err error) {
	t = reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Invalid {
		err = fmt.Errorf("value is invalid:\n %v", s)
		return
	}

	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("provided value is not a struct but %v", t.Kind())
	}

	return
}

// extractTagValues scans given type and tries to find all fields with given tag
// Indices of all found fields are stored as values in resulting map
// Keys of resulting map are actual values of tags
func extractTagValues(typ reflect.Type, tag string) tags {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	tagValues := newTags()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tagValue := field.Tag.Get(tag)

		if field.Anonymous && tagValue != "-" {
			tagValues.Nested[i] = extractTagValues(field.Type, tag)
			continue
		}

		if tagValue != "" && tagValue != "-" && field.PkgPath == "" { // exported
			tagValues.Simple[i] = tagValue
		}

	}

	return tagValues
}

func toMap(obj interface{}, tagmap tags, defaultValue interface{}, policy Policy) (map[string]interface{}, error) {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := map[string]interface{}{}

	for index, tag := range tagmap.Simple {
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

	for index, tags := range tagmap.Nested {
		vField := val.Field(index)
		valueMap, err := toMap(vField.Interface(), tags, defaultValue, policy)
		if err != nil {
			return result, err
		}
		for k, v := range valueMap {
			result[k] = v
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
	case Zeroable:
		if t.IsZero() {
			v = nil
		}
	case ToMappable:
		v, err = t.ToMap()
	}

	return v, nil

}
