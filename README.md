# SToM: Structure To Map converter (yet another one)

## WARNING!
SToM is not supporting nested structs. Stay tuned, it will be fixed soon! :)

## What is it?
Little handy tool to convert your structures into `map[string]interface{}`. It works in 2 modes:
1. General mode. It's when you use exported method `ToMap` to convert arbitrary struct instance to map.
```go
import "github.com/elgris/stom"
import "fmt"

type SomeAwesomeStruct struct {
    ID              int             `db:"id" custom_tag:"id"`
    Name            string          `db:"name"`
    Notes           string
}

func main() {
    s := SomeAwesomeStruct{123, "myname", "mynote"}
    var m map[string]interface{} = stom.ToMap(s)
    fmt.Printf("%+v", m)
    /* you will get map:
        "id": 123,
        "name": "myname"
    Field "Notes" is ignored as it has no tag
    */
}
```

2. "Individual" mode, when you create an instance of SToM for one specific type. In this mode all the tags are analyzed and cached before conversion, thus you can speed the whole process up if you need to convert repeatedly. It's very useful when you need to parse a lot of instances of the same struct.
```go
import "github.com/elgris/stom"
import "fmt"

type SomeAwesomeStruct struct {
    ID              int             `db:"id" custom_tag:"id"`
    Name            string          `db:"name"`
    Notes           string
}

func main() {
    s := SomeAwesomeStruct{123, "myname", "mynote"}
    converter := stom.MustNewStom(s) // at this point 's' is analyzed and tags 'id' and 'name' are cached for future use

    converter.SetTag("db")

    for i:= 0; i < 100500; i++ {
        var m map[string]interface{} = converter.ToMap(s)
    }
}
```

## Benchmarks
https://github.com/elgris/struct-to-map-conversion-benchmark

## TODO
- (???) support filter plugins (???)
- setup travis
- generate godoc and put link to doc

## License
MIT
