# SToM: Structure To Map converter (yet another one)

```go
go get "github.com/elgris/stom"
```

[![GoDoc](https://godoc.org/github.com/elgris/stom?status.png)](https://godoc.org/github.com/elgris/sqrl)
[![Build Status](https://travis-ci.org/elgris/stom.png?branch=master)](https://travis-ci.org/elgris/sqrl)

## What is it?
Little handy tool to convert your structures into `map[string]interface{}`. It works in 2 modes:

**General mode**. It's when you use exported method `ToMap` to convert arbitrary struct instance to map.
```go
import "github.com/elgris/stom"
import "fmt"

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
    ID              int             `db:"id" custom_tag:"id"`
    Name            string          `db:"name"`
    Notes           string
}
func main() {
    s := SomeAwesomeStruct{
        ID:    123,
        Name:  "myname",
        Notes: "mynote",
    }
    s.ParentID = 1123
    s.GrandparentID = 11123

    stom.SetTag("db")
    m, _ :=  stom.ToMap(s)
    fmt.Printf("%+v", m)

    /* you will get map:

        "grand_parent_id": 11123,
        "parent_id": 1123,
        "id": 123,
        "name": "myname"

    Field "Notes" is ignored as it has no tag.
    All embedded structs are flattened into flat map.
    */
}
```

****"Individual" mode**, when you create an instance of SToM for one specific type. In this mode all the tags are analyzed and cached before conversion, thus you can speed the whole process up if you need to convert repeatedly. It's very useful when you need to parse a lot of instances of the same struct.
```go
import "github.com/elgris/stom"
import "fmt"

type SomeAwesomeStruct struct {
    ID              int             `db:"id" custom_tag:"id"`
    Name            string          `db:"name"`
    Notes           string
}

func main() {
    s := SomeAwesomeStruct{
        ID:    123,
        Name:  "myname",
        Notes: "mynote",
    }
    converter := stom.MustNewStom(s) // at this point 's' is analyzed and tags 'id' and 'name' are cached for future use

    converter.SetTag("db")

    for i:= 0; i < 100500; i++ {
        m, _ := converter.ToMap(s)
    }
}
```

## Benchmarks
https://github.com/elgris/struct-to-map-conversion-benchmark

## License
MIT
