package stom

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

type SomeItem struct {
	ID         int             `db:"id" custom_tag:"id"`
	Name       string          `db:"name"`
	Number     int             `db:"number" custom_tag:"num"`
	Checksum   int32           `custom_tag:"sum"`
	Created    time.Time       `custom_tag:"created_time" db:"created"`
	Updated    mysql.NullTime  `db:"updated" custom_tag:"updated_time"`
	Price      float64         `db:"price"`
	Discount   *float64        `db:"discount"`
	IsReserved sql.NullBool    `db:"reserved" custom_tag:"is_reserved"`
	Points     sql.NullInt64   `db:"points"`
	Rating     sql.NullFloat64 `db:"rating"`
	IsVisible  bool            `db:"visible" custom_tag:"visible"`
	Notes      string
}

type testSet struct {
	Item     SomeItem
	Expected map[string]interface{}
}

func TestDefaultTag_DefaultPolicy_DefaultValue(t *testing.T) {
	SetDefault("DEFAULT")
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
			"id":       1,
			"name":     "item_1",
			"number":   11,
			"created":  time.Unix(10000, 0),
			"updated":  "DEFAULT",
			"discount": "DEFAULT",
			"price":    1111.0,
			"reserved": "DEFAULT",
			"points":   "DEFAULT",
			"rating":   "DEFAULT",
			"visible":  true,
		},
	}

	doTest(t, getTestItems(), expecteds)
}

func TestDefaultTag_DefaultPolicy_NilValue(t *testing.T) {
	SetDefault(nil)
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
			"id":       1,
			"name":     "item_1",
			"number":   11,
			"created":  time.Unix(10000, 0),
			"updated":  nil,
			"discount": nil,
			"price":    1111.0,
			"reserved": nil,
			"points":   nil,
			"rating":   nil,
			"visible":  false,
		},
	}

	doTest(t, getTestItems(), expecteds)
}

func TestDefaultTag_ExcludePolicy(t *testing.T) {
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
			"id":      1,
			"name":    "item_1",
			"number":  11,
			"created": time.Unix(10000, 0),
			"price":   1111.0,
			"visible": false,
		},
	}

	doTest(t, getTestItems(), expecteds)

	SetPolicy(PolicyUseDefault)
}

func TestCustomTag_DefaultPolicy_DefaultValue(t *testing.T) {
	SetTags("custom_tag")
	SetDefault("SomeDefault")
	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":           1,
			"num":          11,
			"sum":          111,
			"created_time": time.Unix(10000, 0),
			"updated_time": mysql.NullTime{time.Unix(11000, 0), true},
			"is_reserved":  sql.NullBool{true, true},
			"visible":      true,
		},
		map[string]interface{}{
			"id":           1,
			"num":          11,
			"sum":          111,
			"created_time": "SomeDefault",
			"updated_time": "SomeDefault",
			"is_reserved":  "SomeDefault",
			"visible":      false,
		},
	}

	doTest(t, getTestItems(), expecteds)
}

func TestManyTags_ExcludePolicy_DefaultValue(t *testing.T) {
	SetPolicy(PolicyExclude)
	SetTags("db", "custom_tag")
	SetDefault("SomeDefault")
	expecteds := []map[string]interface{}{
		map[string]interface{}{
			"id":           1,
			"name":         "item_1",
			"number":       11,
			"sum":          111,
			"created_time": time.Unix(10000, 0),
			"updated":      mysql.NullTime{time.Unix(11000, 0), true},
			"discount":     111.0,
			"price":        1111.0,
			"reserved":     sql.NullBool{true, true},
			"points":       sql.NullInt64{int64(11), true},
			"rating":       sql.NullFloat64{1.0, true},
			"visible":      true,
		},
		map[string]interface{}{
			"id":           1,
			"name":         "item_1",
			"number":       11,
			"sum":          111,
			"created_time": time.Unix(10000, 0),
			"price":        1111.0,
			"visible":      false,
		},
	}

	doTest(t, getTestItems(), expecteds)

	SetPolicy(PolicyUseDefault)
}

func doTest(t *testing.T, items []SomeItem, expecteds []map[string]interface{}) {
	if len(items) != len(expecteds) {
		t.Fatalf("number of expected maps %d does not match number of actual items %d",
			len(expecteds),
			len(items))
	}

	for i, set := range items {
		m := ToMap(set.Item)

		for key, expected := range expecteds[i] {
			actual, ok := m[key]
			if !ok {
				t.Fatalf("could not find key %s in map:\n%v\nexpected map:\n%v", key, actual, m)
			}
			assert.Equal(t, v, m[k])
		}
	}
}

func getTestItems() []SomeItem {
	discount := 111.0
	return []SomeItem{
		SomeItem{
			ID:         1,
			Name:       "item_1",
			Number:     11,
			Checksum:   111,
			Created:    time.Unix(10000, 0),
			Updated:    mysql.NullTime{time.Unix(11000, 0), true},
			Price:      1111.0,
			Discount:   &discount,
			IsReserved: sql.NullBool{true, true},
			Points:     sql.NullInt64{int64(11), true},
			Rating:     sql.NullFloat64{1.0, true},
			IsVisible:  true,
			Notes:      "foo",
		},
		SomeItem{
			ID:         1,
			Name:       "item_1",
			Number:     11,
			Checksum:   111,
			Created:    time.Unix(10000, 0),
			Updated:    mysql.NullTime{time.Unix(0, 0), false},
			Price:      1111.0,
			Discount:   nil,
			IsReserved: sql.NullBool{true, false},
			Points:     sql.NullInt64{int64(0), false},
			Rating:     sql.NullFloat64{1.0, false},
			IsVisible:  false,
			Notes:      "bar",
		},
	}
}
