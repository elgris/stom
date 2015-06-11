package stom

import (
	"database/sql"
	"testing"
	"time"

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

type ParentItem struct {
	Base string `db:"base"`
}

type BasicItem struct {
	*ParentItem
	Posted mysql.NullTime `db:"basic_posted"`
}

type AnotherBasicItem struct {
	AnotherBase string `db:"another_base"`
}

type Metainfo struct {
	Tag        string
	Value      string
	MaybeValue sql.NullInt64
	SomeFlag   bool
	Additional map[string]interface{}
}

// ToMap implements ToMappable interface to be used by SToM
func (this Metainfo) ToMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"tag":   this.Tag,
		"value": this.Value,
		"add":   this.Additional,
	}, nil
}

type ComplexItem struct {
	SomeItem
	*BasicItem
	AnotherBasicItem `db:"-"`
	Author           sql.NullString `db:"author"`
	Generation       uint32
	Meta             Metainfo `db:"meta"`
}

func TestComplexItem_DefaultPolicy(t *testing.T) {
	SetTag("db")
	SetDefault("DEFAULT")
	SetPolicy(PolicyUseDefault)

	expected := map[string]interface{}{
		"id":       1,
		"name":     "item_1",
		"number":   11,
		"created":  time.Unix(10000, 0),
		"updated":  mysql.NullTime{time.Unix(11000, 0), true},
		"discount": 111.0,
		"price":    1111.0,
		"reserved": sql.NullBool{true, true},
		"points":   "DEFAULT",
		"rating":   sql.NullFloat64{1.0, true},
		"visible":  true,

		"base":         "base",
		"basic_posted": "DEFAULT",

		"author": "DEFAULT",
		"meta": map[string]interface{}{
			"tag":   "metatag",
			"value": "valvalval",
			"add": map[string]interface{}{
				"foo": 12,
				"bar": sql.NullBool{true, false},
			},
		},
	}

	doTest(t, getTestComplexItem(), expected)
}

func TestComplexItem_ExcludePolicy(t *testing.T) {
	SetTag("db")
	SetDefault("DEFAULT")
	SetPolicy(PolicyExclude)

	expected := map[string]interface{}{
		"id":       1,
		"name":     "item_1",
		"number":   11,
		"created":  time.Unix(10000, 0),
		"updated":  mysql.NullTime{time.Unix(11000, 0), true},
		"discount": 111.0,
		"price":    1111.0,
		"reserved": sql.NullBool{true, true},
		"rating":   sql.NullFloat64{1.0, true},
		"visible":  true,

		"base": "base",

		"meta": map[string]interface{}{
			"tag":   "metatag",
			"value": "valvalval",
			"add": map[string]interface{}{
				"foo": 12,
				"bar": sql.NullBool{true, false},
			},
		},
	}

	doTest(t, getTestComplexItem(), expected)
}

func TestDefaultPolicy_DefaultValue(t *testing.T) {
	SetTag("db")
	SetDefault("DEFAULT")
	SetPolicy(PolicyUseDefault)

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

	doTestItems(t, getTestItems(), expecteds)
}

func TestDefaultPolicy_NilValue(t *testing.T) {
	SetTag("db")
	SetDefault(nil)
	SetPolicy(PolicyUseDefault)

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

	doTestItems(t, getTestItems(), expecteds)
}

func TestExcludePolicy(t *testing.T) {
	SetTag("db")
	SetPolicy(PolicyExclude)
	SetDefault("SomeDefault")

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

	doTestItems(t, getTestItems(), expecteds)
}

func TestCustomTag_DefaultPolicy(t *testing.T) {
	SetTag("custom_tag")
	SetDefault("SomeDefault")
	SetPolicy(PolicyUseDefault)

	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":               1,
			"num":              11,
			"sum":              111,
			"created_time":     time.Unix(10000, 0),
			"updated_time":     mysql.NullTime{time.Unix(11000, 0), true},
			"is_reserved":      sql.NullBool{true, true},
			"visible":          true,
			"i_ignore_nothing": 10,
		},
		map[string]interface{}{
			"id":               2,
			"num":              22,
			"sum":              222,
			"created_time":     time.Unix(20000, 0),
			"updated_time":     "SomeDefault",
			"is_reserved":      "SomeDefault",
			"visible":          false,
			"i_ignore_nothing": 20,
		},
	}

	doTestItems(t, getTestItems(), expecteds)
}

func doTestItems(t *testing.T, items []SomeItem, expecteds []map[string]interface{}) {
	if len(items) != len(expecteds) {
		t.Fatalf("number of expected maps %d does not match number of actual items %d",
			len(expecteds),
			len(items))
	}

	for i := range items {
		doTest(t, items[i], expecteds[i])
	}
}

func doTest(t *testing.T, item interface{}, expected map[string]interface{}) {
	actual, err := ConvertToMap(item)
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

func getTestComplexItem() ComplexItem {
	discount := 111.0
	item := ComplexItem{
		SomeItem: SomeItem{
			ID:              1,
			Name:            "item_1",
			Number:          11,
			Checksum:        111,
			Created:         time.Unix(10000, 0),
			Updated:         mysql.NullTime{time.Unix(11000, 0), true},
			Price:           1111.0,
			Discount:        &discount,
			IsReserved:      sql.NullBool{true, true},
			Points:          sql.NullInt64{int64(11), false},
			Rating:          sql.NullFloat64{1.0, true},
			IsVisible:       true,
			Notes:           "foo",
			SomeIgnoreField: 10,
		},
		BasicItem: &BasicItem{
			ParentItem: &ParentItem{Base: "base"},
			Posted:     mysql.NullTime{time.Now(), false},
		},
		AnotherBasicItem: AnotherBasicItem{
			AnotherBase: "anotherBase",
		},
		Author:     sql.NullString{"invalid_author", false},
		Generation: 123,
	}

	item.Meta.Tag = "metatag"
	item.Meta.Value = "valvalval"
	item.Meta.MaybeValue = sql.NullInt64{int64(11), true}
	item.Meta.SomeFlag = true
	item.Meta.Additional = map[string]interface{}{
		"foo": 12,
		"bar": sql.NullBool{true, false},
	}

	return item
}
