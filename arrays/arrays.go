package arrays

func Sum(nums [5]int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}
