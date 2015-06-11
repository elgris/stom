package stom

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

type ComplexItem struct {
	SomeItem
	*BasicItem
	AnotherBasicItem `db:"-"`
	Author           sql.NullString `db:"author"`
	Generation       uint32
	Meta             Metainfo `db:"meta"`
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
