package arrays

func Sum(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func SumAll(numSlices ...[]int) []int {
	var sums []int
	for _, nums := range numSlices {
		sums = append(sums, Sum(nums))
	}
	return sums
}
