package stom_test

import (
	"fmt"

	"github.com/elgris/stom"
)

type SomeGrandparentStruct struct {
	GrandparentID int `db:"grand_parent_id"`
}

type SomeParentStruct struct {
	SomeGrandparentStruct
	ParentID int `db:"parent_id"`
}

// Yes, SToM supports embedded structs.
type SomeAwesomeStruct struct {
	SomeParentStruct
	ID            int         `db:"id" custom_tag:"id"`
	Name          string      `db:"name"`
	AbstractThing interface{} `db:"thing"`
	Notes         string
}

func Example() {
	s := SomeAwesomeStruct{
		ID:    123,
		Name:  "myname",
		Notes: "mynote",
	}
	s.ParentID = 1123
	s.GrandparentID = 11123

	converter := stom.MustNewStom(s).
		SetTag("db").
		SetPolicy(stom.PolicyExclude).
		SetDefault("DEFAULT")

	/* you will get map:

	       "grand_parent_id": 11123,
	       "parent_id": 1123,
	       "id": 123,
	       "name": "myname",
	       "thing": "DEFAULT"

	   Field "Notes" is ignored as it has no tag.
	   Field "AbstractThing" is nil, so it replaced with default value.
	   Current policy demands to use default instead of nil values.

	   All embedded structs are flattened into flat map.
	*/
	m, err := converter.ToMap(s)
	fmt.Printf("MAP: %+v\nERROR: %v", m, err)
}
