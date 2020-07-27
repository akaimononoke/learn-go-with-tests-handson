package iteration

import (
	"fmt"
	"strings"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 5)
	expected := "aaaaa"
	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}

func ExampleRepeat() {
	repeated := Repeat("b", 10)
	fmt.Println(repeated)
	// Output: bbbbbbbbbb
}

func BenchmarkCompare_Operator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = "a" < "b"
	}
}

func BenchmarkCompare_StringsCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.Compare("a", "b")
	}
}
