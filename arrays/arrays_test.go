package arrays

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("collection of 5 numbers", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		got := Sum(nums)
		want := 15
		if got != want {
			t.Errorf("with %v, got %d, want %d", nums, got, want)
		}
	})
}

func BenchmarkSum(b *testing.B) {
	nums := []int{0, 2, 4, 6, 8, 10}
	for i := 0; i < b.N; i++ {
		Sum(nums)
	}
}

func ExampleSum() {
	nums := []int{0, 3, 6}
	sum := Sum(nums)
	fmt.Println(sum)
	// Output: 9
}
