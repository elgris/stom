package stom

import (
	"testing"
)

func BenchmarkDefaultPolicy_DefaultValue0(b *testing.B) {
	SetPolicy(PolicyUseDefault)
	SetTag("db")
	SetDefault("SomeDefault")

	items := getTestItems()
	doBenchmark(b, items[0])
}
func BenchmarkDefaultPolicy_DefaultValue1(b *testing.B) {
	SetPolicy(PolicyUseDefault)
	SetTag("db")
	SetDefault("SomeDefault")

	items := getTestItems()
	doBenchmark(b, items[1])
}

func BenchmarkExcludePolicy0(b *testing.B) {
	SetPolicy(PolicyExclude)
	SetTag("db")
	SetDefault("SomeDefault")

	items := getTestItems()
	doBenchmark(b, items[0])
}
func BenchmarkExcludePolicy1(b *testing.B) {
	SetPolicy(PolicyExclude)
	SetTag("db")
	SetDefault("SomeDefault")

	items := getTestItems()
	doBenchmark(b, items[1])
}

func doBenchmark(b *testing.B, item SomeItem) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToMap(item)
	}
}
