package stom_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/elgris/stom"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

type SomeItem struct {
	ID              int             `db:"id" custom_tag:"id"`
	Name            string          `db:"name"`
	somePrivate     string          `db:"some_private" custom_tag:"some_private"`
	Number          int             `db:"number" custom_tag:"num"`
	Checksum        int32           `custom_tag:"sum"`
	Created         time.Time       `custom_tag:"created_time" db:"created"`
	Updated         mysql.NullTime  `db:"updated" custom_tag:"updated_time"`
	Price           float64         `db:"price"`
	Discount        *float64        `db:"discount"`
	IsReserved      sql.NullBool    `db:"reserved" custom_tag:"is_reserved"`
	Points          sql.NullInt64   `db:"points"`
	Rating          sql.NullFloat64 `db:"rating"`
	IsVisible       bool            `db:"visible" custom_tag:"visible"`
	SomeIgnoreField int             `db:"-" custom_tag:"i_ignore_nothing"`
	Notes           string
}

func TestDefaultPolicy_DefaultValue(t *testing.T) {
	stom.SetTag("db")
	stom.SetDefault("DEFAULT")
	stom.SetPolicy(stom.PolicyUseDefault)

	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":       1,
			"name":     "item_1",
			"number":   11,
			"created":  time.Unix(10000, 0),
			"updated":  mysql.NullTime{time.Unix(11000, 0), true},
			"discount": 111.0,
			"price":    1111.0,
			"reserved": sql.NullBool{true, true},
			"points":   sql.NullInt64{int64(11), true},
			"rating":   sql.NullFloat64{1.0, true},
			"visible":  true,
		},
		map[string]interface{}{
			"id":       2,
			"name":     "item_2",
			"number":   22,
			"created":  time.Unix(20000, 0),
			"updated":  "DEFAULT",
			"discount": "DEFAULT",
			"price":    2222.0,
			"reserved": "DEFAULT",
			"points":   "DEFAULT",
			"rating":   "DEFAULT",
			"visible":  false,
		},
	}

	doTestItems(t, stom.ToMapperFunc(stom.ConvertToMap), getTestItems(), expecteds)
}

func TestDefaultPolicy_NilValue(t *testing.T) {
	stom.SetTag("db")
	stom.SetDefault(nil)
	stom.SetPolicy(stom.PolicyUseDefault)

	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":       1,
			"name":     "item_1",
			"number":   11,
			"created":  time.Unix(10000, 0),
			"updated":  mysql.NullTime{time.Unix(11000, 0), true},
			"discount": 111.0,
			"price":    1111.0,
			"reserved": sql.NullBool{true, true},
			"points":   sql.NullInt64{int64(11), true},
			"rating":   sql.NullFloat64{1.0, true},
			"visible":  true,
		},
		map[string]interface{}{
			"id":       2,
			"name":     "item_2",
			"number":   22,
			"created":  time.Unix(20000, 0),
			"updated":  nil,
			"discount": nil,
			"price":    2222.0,
			"reserved": nil,
			"points":   nil,
			"rating":   nil,
			"visible":  false,
		},
	}

	doTestItems(t, stom.ToMapperFunc(stom.ConvertToMap), getTestItems(), expecteds)
}

func TestExcludePolicy(t *testing.T) {
	stom.SetTag("db")
	stom.SetPolicy(stom.PolicyExclude)
	stom.SetDefault("SomeDefault")

	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":       1,
			"name":     "item_1",
			"number":   11,
			"created":  time.Unix(10000, 0),
			"updated":  mysql.NullTime{time.Unix(11000, 0), true},
			"discount": 111.0,
			"price":    1111.0,
			"reserved": sql.NullBool{true, true},
			"points":   sql.NullInt64{int64(11), true},
			"rating":   sql.NullFloat64{1.0, true},
			"visible":  true,
		},
		map[string]interface{}{
			"id":      2,
			"name":    "item_2",
			"number":  22,
			"created": time.Unix(20000, 0),
			"price":   2222.0,
			"visible": false,
		},
	}

	doTestItems(t, stom.ToMapperFunc(stom.ConvertToMap), getTestItems(), expecteds)
}

func TestToMapper_CustomTag_DefaultPolicy(t *testing.T) {
	stom.SetTag("db")
	stom.SetPolicy(stom.PolicyExclude)
	stom.SetDefault("foo")

	tomapper := stom.MustNewStom(SomeItem{}).
		SetTag("custom_tag").
		SetPolicy(stom.PolicyUseDefault).
		SetDefault("SomeDefault")

	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":               1,
			"num":              11,
			"sum":              int32(111),
			"created_time":     time.Unix(10000, 0),
			"updated_time":     mysql.NullTime{time.Unix(11000, 0), true},
			"is_reserved":      sql.NullBool{true, true},
			"visible":          true,
			"i_ignore_nothing": 10,
		},
		map[string]interface{}{
			"id":               2,
			"num":              22,
			"sum":              int32(222),
			"created_time":     time.Unix(20000, 0),
			"updated_time":     "SomeDefault",
			"is_reserved":      "SomeDefault",
			"visible":          false,
			"i_ignore_nothing": 20,
		},
	}

	doTestItems(t, tomapper, getTestItems(), expecteds)
}

func TestCustomTag_DefaultPolicy(t *testing.T) {
	stom.SetTag("custom_tag")
	stom.SetDefault("SomeDefault")
	stom.SetPolicy(stom.PolicyUseDefault)

	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":               1,
			"num":              11,
			"sum":              int32(111),
			"created_time":     time.Unix(10000, 0),
			"updated_time":     mysql.NullTime{time.Unix(11000, 0), true},
			"is_reserved":      sql.NullBool{true, true},
			"visible":          true,
			"i_ignore_nothing": 10,
		},
		map[string]interface{}{
			"id":               2,
			"num":              22,
			"sum":              int32(222),
			"created_time":     time.Unix(20000, 0),
			"updated_time":     "SomeDefault",
			"is_reserved":      "SomeDefault",
			"visible":          false,
			"i_ignore_nothing": 20,
		},
	}

	doTestItems(t, stom.ToMapperFunc(stom.ConvertToMap), getTestItems(), expecteds)
}

func doTestItems(t *testing.T, tomapper stom.ToMapper, items []SomeItem, expecteds []map[string]interface{}) {
	if len(items) != len(expecteds) {
		t.Fatalf("number of expected maps %d does not match number of actual items %d",
			len(expecteds),
			len(items))
	}

	for i := range items {
		doTest(t, tomapper, items[i], expecteds[i])
	}
}

func doTest(t *testing.T, tomapper stom.ToMapper, item interface{}, expected map[string]interface{}) {
	actual, err := tomapper.ToMap(item)
	if err != nil {
		t.Fatalf("ToMap call returned error: %s", err.Error())
	}

	// TODO: maybe use just equal asserion here?
	if len(actual) != len(expected) {
		t.Fatalf("size of expected map %d\n%+v\ndoes not match size of generated map %d\n%+v",
			len(expected),
			expected,
			len(actual),
			actual)
	}

	for key, e := range expected {
		a, ok := actual[key]
		if !ok {
			t.Fatalf("could not find key %s in map:\n%v\nexpected map:\n%v", key, actual, expected)
		}
		if !assert.Equal(t, e, a) {
			t.Fatalf("expected value by key %s is:\n %#v\ngot:\n %#v", key, e, a)
		}
	}
}

func getTestItems() []SomeItem {
	discount := 111.0
	return []SomeItem{
		SomeItem{
			ID:              1,
			Name:            "item_1",
			Number:          11,
			Checksum:        111,
			Created:         time.Unix(10000, 0),
			Updated:         mysql.NullTime{time.Unix(11000, 0), true},
			Price:           1111.0,
			Discount:        &discount,
			IsReserved:      sql.NullBool{true, true},
			Points:          sql.NullInt64{int64(11), true},
			Rating:          sql.NullFloat64{1.0, true},
			IsVisible:       true,
			Notes:           "foo",
			SomeIgnoreField: 10,
		},
		SomeItem{
			ID:              2,
			Name:            "item_2",
			Number:          22,
			Checksum:        222,
			Created:         time.Unix(20000, 0),
			Updated:         mysql.NullTime{time.Unix(0, 0), false},
			Price:           2222.0,
			Discount:        nil,
			IsReserved:      sql.NullBool{true, false},
			Points:          sql.NullInt64{int64(0), false},
			Rating:          sql.NullFloat64{2.0, false},
			IsVisible:       false,
			Notes:           "bar",
			SomeIgnoreField: 20,
		},
	}
}
